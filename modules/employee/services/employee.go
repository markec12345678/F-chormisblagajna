package services

import (
	"context"
	"sort"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules/core/models"
	employee_models "github.com/nutrixpos/pos/modules/employee/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type EmployeeService struct {
	Logger logger.ILogger
	Config config.Config
}

func (s *EmployeeService) GetPerformance(startDate, endDate time.Time) (*employee_models.PerformanceSummary, error) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ordersCollection := client.Database(s.Config.Databases[0].Database).Collection("orders")

	filter := bson.M{
		"submitted_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"is_paid": true,
	}

	cursor, err := ordersCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	employeeMap := make(map[string]*employee_models.EmployeePerformance)

	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			continue
		}

		employeeId := order.CustomData["employee_id"]
		employeeName := order.CustomData["employee_name"]
		if employeeId == "" {
			employeeId = "unknown"
			employeeName = "Unknown"
		}

		if _, exists := employeeMap[employeeId]; !exists {
			employeeMap[employeeId] = &employee_models.EmployeePerformance{
				EmployeeId:   employeeId,
				EmployeeName: employeeName,
				TopProducts:  make([]employee_models.ProductStat, 0),
			}
		}

		perf := employeeMap[employeeId]
		perf.OrderCount++
		saleAmount := order.SalePrice - order.Discount
		perf.TotalSales += saleAmount
		perf.TotalTips += order.Tips

		productMap := make(map[string]*employee_models.ProductStat)
		for _, p := range perf.TopProducts {
			productMap[p.ProductId] = &p
		}

		for _, item := range order.Items {
			pid := item.Product.Id
			if _, exists := productMap[pid]; !exists {
				productMap[pid] = &employee_models.ProductStat{
					ProductId:   pid,
					ProductName: item.Product.Name,
				}
			}
			productMap[pid].Quantity += int(item.Quantity)
			productMap[pid].Revenue += item.Price * item.Quantity
		}

		perf.TopProducts = make([]employee_models.ProductStat, 0, len(productMap))
		for _, p := range productMap {
			perf.TopProducts = append(perf.TopProducts, *p)
		}
	}

	employees := make([]employee_models.EmployeePerformance, 0, len(employeeMap))
	for _, emp := range employeeMap {
		if emp.OrderCount > 0 {
			emp.AvgOrderValue = emp.TotalSales / float64(emp.OrderCount)
		}
		emp.SalesPerHour = emp.TotalSales / 8.0
		employees = append(employees, *emp)
	}

	sort.Slice(employees, func(i, j int) bool {
		return employees[i].TotalSales > employees[j].TotalSales
	})
	for i := range employees {
		employees[i].Rank = i + 1
	}

	var totalRevenue float64
	for _, emp := range employees {
		totalRevenue += emp.TotalSales
	}

	var totalSalesPerHour float64
	for _, emp := range employees {
		totalSalesPerHour += emp.SalesPerHour
	}
	avgSalesPerHour := 0.0
	if len(employees) > 0 {
		avgSalesPerHour = totalSalesPerHour / float64(len(employees))
	}

	summary := &employee_models.PerformanceSummary{
		TotalEmployees:  len(employees),
		TopPerformers:   employees,
		TotalRevenue:    totalRevenue,
		AvgSalesPerHour: avgSalesPerHour,
	}

	return summary, nil
}
