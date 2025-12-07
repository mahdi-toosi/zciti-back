# Add Database Relation

Add a new relation between existing models.

## Relation Types

### One-to-Many (e.g., User has many Orders)
```go
// In User schema
Orders []Order `gorm:"foreignKey:UserID"`

// In Order schema
UserID uint64 `json:",omitempty"`
User   *User  `gorm:"foreignKey:UserID" json:",omitempty"`
```

### Many-to-Many (e.g., User belongs to many Businesses)
```go
// In User schema
Businesses []Business `gorm:"many2many:business_users;"`

// In Business schema
Users []User `gorm:"many2many:business_users;"`
```

### One-to-One (e.g., User has one Profile)
```go
// In User schema
Profile *Profile `gorm:"foreignKey:UserID"`

// In Profile schema
UserID uint64 `gorm:"unique"`
User   *User  `gorm:"foreignKey:UserID"`
```

## Example Usage

"Add a many-to-many relation between Product and Category models"

