package faker

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"gorm.io/gorm"
)

func UserFaker(db *gorm.DB) []*entities.UserModels {
	customers := []*entities.UserModels{
		{
			ID:           2,
			Email:        "admin@gmail.com",
			Password:     "Q9yxsgcFqch0DUYAxTq++2EDCgXP3Why1fwkRWe20HFEtbZYNp1/0AIvymKV8sVA",
			Phone:        "081234567891",
			Role:         "admin",
			Name:         "Admin Disappear",
			IsVerified:   true,
			PhotoProfile: "https://res.cloudinary.com/dufa4bel6/image/upload/v1702738826/disappear/NEW_LOGO_qamyiz.jpg",
		},
		{
			ID:           3,
			Email:        "customer2@example.com",
			Password:     "aq9CgZMSuxopX0iezFvLYIwmJgD9TlZ4PReX2zgUK0Nxh3TUTYvzGYZt8PAL505A",
			Phone:        "081234567892",
			Role:         "customer",
			Name:         "Mijan",
			IsVerified:   true,
			PhotoProfile: "https://res.cloudinary.com/dufa4bel6/image/upload/v1702299949/disappear/isxjm4ccc6rh2n4rwwik.jpg",
		},
		{
			ID:           4,
			Email:        "customer3@example.com",
			Password:     "aq9CgZMSuxopX0iezFvLYIwmJgD9TlZ4PReX2zgUK0Nxh3TUTYvzGYZt8PAL505A",
			Phone:        "081234567893",
			Role:         "customer",
			Name:         "Samini",
			IsVerified:   true,
			PhotoProfile: "https://res.cloudinary.com/dufa4bel6/image/upload/v1702299949/disappear/isxjm4ccc6rh2n4rwwik.jpg",
		},
		{
			ID:           5,
			Email:        "customer4@example.com",
			Password:     "aq9CgZMSuxopX0iezFvLYIwmJgD9TlZ4PReX2zgUK0Nxh3TUTYvzGYZt8PAL505A",
			Phone:        "081234567894",
			Role:         "customer",
			Name:         "Joko",
			IsVerified:   true,
			PhotoProfile: "https://res.cloudinary.com/dufa4bel6/image/upload/v1702299949/disappear/isxjm4ccc6rh2n4rwwik.jpg",
		},
		{
			ID:           6,
			Email:        "customer5@example.com",
			Password:     "aq9CgZMSuxopX0iezFvLYIwmJgD9TlZ4PReX2zgUK0Nxh3TUTYvzGYZt8PAL505A",
			Phone:        "081234567894",
			Role:         "customer",
			Name:         "Siti Jubaidah",
			IsVerified:   true,
			PhotoProfile: "https://res.cloudinary.com/dufa4bel6/image/upload/v1702299949/disappear/isxjm4ccc6rh2n4rwwik.jpg",
		},
	}
	return customers
}
