package dao

import (
	"PersonalityTest/internal/dao/internal"
)

// ordersDao is the data access object for the table orders.
type ordersDao struct {
	*internal.OrdersDao
}

var (
	// Orders is a globally accessible object for table orders operations.
	Orders = ordersDao{internal.NewOrdersDao()}
)

// Add your custom methods and functionality below.
