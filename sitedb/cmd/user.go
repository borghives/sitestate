package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Define the "user" context command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage MongoDB users",
}

// Define the "create" action command
var createCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a MongoDB user",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")

		if name == "" || password == "" {
			log.Fatalf("Name and password are required")
		}

		readDb, _ := cmd.Flags().GetStringSlice("read")
		readWriteDb, _ := cmd.Flags().GetStringSlice("write")

		fmt.Printf("Action: Set MongoDB user '%s'...\n", name)
		fmt.Printf("Read permission: %v\n", readDb)
		fmt.Printf("ReadWrite permission: %v\n", readWriteDb)
		client := GetDbClient(cmd)
		defer client.Disconnect(context.Background())

		err := UpsertDbUser(client, name, password, readDb, readWriteDb, false)
		if err != nil {
			log.Fatalf("Failed to set admin: %v", err)
		}
	},
}

func translateRole(readDb []string, readWriteDb []string, isAdmin bool) bson.A {
	roles := bson.A{}
	for _, db := range readDb {
		roles = append(roles, bson.M{"role": "read", "db": db})
	}
	for _, db := range readWriteDb {
		roles = append(roles, bson.M{"role": "readWrite", "db": db})
	}
	if isAdmin {
		roles = append(roles, bson.M{"role": "userAdminAnyDatabase", "db": "admin"})
	}
	return roles
}

func CreateDbUser(client *mongo.Client, username string, newPassword string, readDb []string, readWriteDb []string, isAdmin bool) error {
	createUserCmd := bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: newPassword},
		{Key: "roles", Value: translateRole(readDb, readWriteDb, isAdmin)},
	}

	var result bson.M
	err := client.Database("admin").RunCommand(context.Background(), createUserCmd).Decode(&result)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created user: %s\n", username)
	return nil
}

func UpdateDbUser(client *mongo.Client, username string, newPassword string, readDb []string, readWriteDb []string, isAdmin bool) error {

	roles := translateRole(readDb, readWriteDb, isAdmin)
	fmt.Printf("Roles: %v\n", roles)
	updateUserCmd := bson.D{
		{Key: "updateUser", Value: username},
		{Key: "pwd", Value: newPassword},
		{Key: "roles", Value: roles},
	}

	var result bson.M
	err := client.Database("admin").RunCommand(context.Background(), updateUserCmd).Decode(&result)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated user: %s\n", username)
	return nil
}

func UpsertDbUser(client *mongo.Client, username string, newPassword string, readDb []string, readWriteDb []string, isAdmin bool) error {
	err := UpdateDbUser(client, username, newPassword, readDb, readWriteDb, isAdmin)
	if err != nil {
		return CreateDbUser(client, username, newPassword, readDb, readWriteDb, isAdmin)
	}
	return nil
}

var readDb []string
var readWriteDb []string

func init() {
	// Add the action to the context
	userCmd.AddCommand(createCmd)

	// Add the context to the root dbenv command
	rootCmd.AddCommand(userCmd)

	// Define flags specifically for the 'create' action
	createCmd.Flags().StringP("name", "n", "", "Username for the new user")
	createCmd.Flags().StringP("password", "p", "", "Password for the new user")
	createCmd.Flags().StringSliceVarP(&readDb, "read", "r", []string{}, "List of read database (comma-separated or multiple flags)")
	createCmd.Flags().StringSliceVarP(&readWriteDb, "write", "w", []string{}, "List of readWrite database (comma-separated or multiple flags)")
}
