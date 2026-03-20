package data

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UpdateOperator struct {
	currentDate bson.M
	inc         bson.M
	min         bson.M
	max         bson.M
	mul         bson.M
	rename      bson.M
	set         bson.M
	setDoc      interface{}
	setOnInsert bson.M
	unset       bson.M
}

func NewUpdate() *UpdateOperator {
	return &UpdateOperator{}
}

func (u *UpdateOperator) CurrentDate(key string) *UpdateOperator {
	if u.currentDate == nil {
		u.currentDate = bson.M{}
	}
	u.currentDate[key] = bson.M{"$type": "date"}
	return u
}

func (u *UpdateOperator) Inc(key string, value interface{}) *UpdateOperator {
	if u.inc == nil {
		u.inc = bson.M{}
	}
	u.inc[key] = value
	return u
}

func (u *UpdateOperator) Min(key string, value interface{}) *UpdateOperator {
	if u.min == nil {
		u.min = bson.M{}
	}
	u.min[key] = value
	return u
}

func (u *UpdateOperator) Max(key string, value interface{}) *UpdateOperator {
	if u.max == nil {
		u.max = bson.M{}
	}
	u.max[key] = value
	return u
}

func (u *UpdateOperator) Mul(key string, value interface{}) *UpdateOperator {
	if u.mul == nil {
		u.mul = bson.M{}
	}
	u.mul[key] = value
	return u
}

func (u *UpdateOperator) Rename(key string, value interface{}) *UpdateOperator {
	if u.rename == nil {
		u.rename = bson.M{}
	}
	u.rename[key] = value
	return u
}

func (u *UpdateOperator) Set(key string, value interface{}) *UpdateOperator {

	if u.setDoc != nil {
		log.Fatal("Set by document in used.  Only one type of set can be used.  Either set by field name Set() or by document SetDoc()")
		return u
	}

	if u.set == nil {
		u.set = bson.M{}
	}
	u.set[key] = value
	return u
}

func (u *UpdateOperator) SetDoc(doc interface{}) *UpdateOperator {
	if u.set != nil {
		log.Fatal("Set by key field in used.  Only one type of set can be used.  Either set by field name Set() or by document SetDoc()")
		return u
	}
	u.setDoc = doc
	return u
}

func (u *UpdateOperator) SetOnInsert(key string, value interface{}) *UpdateOperator {
	if u.setOnInsert == nil {
		u.setOnInsert = bson.M{}
	}
	u.setOnInsert[key] = value
	return u
}

func (u *UpdateOperator) Unset(key string) *UpdateOperator {
	if u.unset == nil {
		u.unset = bson.M{}
	}
	u.unset[key] = ""
	return u
}

func (u *UpdateOperator) ToPrimitive() bson.M {
	retval := bson.M{}

	if u.set != nil {
		retval["$set"] = u.set
	} else if u.setDoc != nil {
		retval["$set"] = u.setDoc
	}

	if u.unset != nil {
		retval["$unset"] = u.unset
	}
	if u.inc != nil {
		retval["$inc"] = u.inc
	}
	if u.min != nil {
		retval["$min"] = u.min
	}
	if u.max != nil {
		retval["$max"] = u.max
	}
	if u.mul != nil {
		retval["$mul"] = u.mul
	}
	if u.rename != nil {
		retval["$rename"] = u.rename
	}
	if u.setOnInsert != nil {
		retval["$setOnInsert"] = u.setOnInsert
	}
	if u.currentDate != nil {
		retval["$currentDate"] = u.currentDate
	}

	return retval
}
