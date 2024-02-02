package data

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateOperator struct {
	currentDate primitive.M
	inc         primitive.M
	min         primitive.M
	max         primitive.M
	mul         primitive.M
	rename      primitive.M
	set         primitive.M
	setDoc      interface{}
	setOnInsert primitive.M
	unset       primitive.M
}

func NewUpdate() *UpdateOperator {
	return &UpdateOperator{}
}

func (u *UpdateOperator) CurrentDate(key string) *UpdateOperator {
	if u.currentDate == nil {
		u.currentDate = primitive.M{}
	}
	u.currentDate[key] = primitive.M{"$type": "date"}
	return u
}

func (u *UpdateOperator) Inc(key string, value interface{}) *UpdateOperator {
	if u.inc == nil {
		u.inc = primitive.M{}
	}
	u.inc[key] = value
	return u
}

func (u *UpdateOperator) Min(key string, value interface{}) *UpdateOperator {
	if u.min == nil {
		u.min = primitive.M{}
	}
	u.min[key] = value
	return u
}

func (u *UpdateOperator) Max(key string, value interface{}) *UpdateOperator {
	if u.max == nil {
		u.max = primitive.M{}
	}
	u.max[key] = value
	return u
}

func (u *UpdateOperator) Mul(key string, value interface{}) *UpdateOperator {
	if u.mul == nil {
		u.mul = primitive.M{}
	}
	u.mul[key] = value
	return u
}

func (u *UpdateOperator) Rename(key string, value interface{}) *UpdateOperator {
	if u.rename == nil {
		u.rename = primitive.M{}
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
		u.set = primitive.M{}
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
		u.setOnInsert = primitive.M{}
	}
	u.setOnInsert[key] = value
	return u
}

func (u *UpdateOperator) Unset(key string) *UpdateOperator {
	if u.unset == nil {
		u.unset = primitive.M{}
	}
	u.unset[key] = ""
	return u
}

func (u *UpdateOperator) ToPrimitive() primitive.M {
	retval := primitive.M{}

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
