// Package modules_test verifies that all module types implement required interfaces.
package modules_test

import (
	"github.com/nutrixpos/pos/modules"
	"github.com/nutrixpos/pos/modules/accounting"
	"github.com/nutrixpos/pos/modules/ai"
	"github.com/nutrixpos/pos/modules/auditlog"
	"github.com/nutrixpos/pos/modules/branch"
	"github.com/nutrixpos/pos/modules/chat"
	"github.com/nutrixpos/pos/modules/core"
	"github.com/nutrixpos/pos/modules/customerdisplay"
	"github.com/nutrixpos/pos/modules/delivery"
	"github.com/nutrixpos/pos/modules/employee"
	"github.com/nutrixpos/pos/modules/expense"
	"github.com/nutrixpos/pos/modules/feedback"
	"github.com/nutrixpos/pos/modules/fiscal"
	fiscal_hr "github.com/nutrixpos/pos/modules/fiscal_hr"
	"github.com/nutrixpos/pos/modules/floorplan"
	"github.com/nutrixpos/pos/modules/giftcards"
	"github.com/nutrixpos/pos/modules/hubsync"
	"github.com/nutrixpos/pos/modules/inventoryalerts"
	"github.com/nutrixpos/pos/modules/inventorytransfer"
	"github.com/nutrixpos/pos/modules/kiosk"
	"github.com/nutrixpos/pos/modules/kitchen"
	"github.com/nutrixpos/pos/modules/loyalty"
	"github.com/nutrixpos/pos/modules/marketing"
	"github.com/nutrixpos/pos/modules/menuengineering"
	"github.com/nutrixpos/pos/modules/multilocation"
	"github.com/nutrixpos/pos/modules/multipayment"
	"github.com/nutrixpos/pos/modules/notification"
	"github.com/nutrixpos/pos/modules/onlineorder"
	"github.com/nutrixpos/pos/modules/promotion"
	"github.com/nutrixpos/pos/modules/purchase"
	"github.com/nutrixpos/pos/modules/queue"
	"github.com/nutrixpos/pos/modules/receipt"
	"github.com/nutrixpos/pos/modules/report"
	"github.com/nutrixpos/pos/modules/reservations"
	"github.com/nutrixpos/pos/modules/scheduling"
	"github.com/nutrixpos/pos/modules/splitbill"
	"github.com/nutrixpos/pos/modules/supplier"
	"github.com/nutrixpos/pos/modules/table"
	"github.com/nutrixpos/pos/modules/tableside"
	"github.com/nutrixpos/pos/modules/timeclock"
	"github.com/nutrixpos/pos/modules/tips"
	"github.com/nutrixpos/pos/modules/training"
	"github.com/nutrixpos/pos/modules/waste"
)

var _ modules.IHttpModule = (*accounting.AccountingModule)(nil)
var _ modules.IHttpModule = (*ai.AIModule)(nil)
var _ modules.IHttpModule = (*auditlog.AuditLogModule)(nil)
var _ modules.IHttpModule = (*branch.BranchModule)(nil)
var _ modules.IHttpModule = (*chat.ChatModule)(nil)
var _ modules.IHttpModule = (*core.Core)(nil)
var _ modules.IHttpModule = (*customerdisplay.CustomerDisplayModule)(nil)
var _ modules.IHttpModule = (*delivery.DeliveryModule)(nil)
var _ modules.IHttpModule = (*employee.EmployeeModule)(nil)
var _ modules.IHttpModule = (*expense.ExpenseModule)(nil)
var _ modules.IHttpModule = (*feedback.FeedbackModule)(nil)
var _ modules.IHttpModule = (*fiscal.FiscalModule)(nil)
var _ modules.IHttpModule = (*fiscal_hr.FiscalModuleHR)(nil)
var _ modules.IHttpModule = (*floorplan.FloorplanModule)(nil)
var _ modules.IHttpModule = (*giftcards.GiftCardModule)(nil)
var _ modules.IHttpModule = (*hubsync.HubSyncModule)(nil)
var _ modules.IHttpModule = (*inventoryalerts.InventoryAlertsModule)(nil)
var _ modules.IHttpModule = (*inventorytransfer.InventoryTransferModule)(nil)
var _ modules.IHttpModule = (*kiosk.KioskModule)(nil)
var _ modules.IHttpModule = (*kitchen.KitchenModule)(nil)
var _ modules.IHttpModule = (*loyalty.LoyaltyModule)(nil)
var _ modules.IHttpModule = (*marketing.MarketingModule)(nil)
var _ modules.IHttpModule = (*menuengineering.MenuEngineeringModule)(nil)
var _ modules.IHttpModule = (*multilocation.MultiLocationModule)(nil)
var _ modules.IHttpModule = (*multipayment.MultiPaymentModule)(nil)
var _ modules.IHttpModule = (*notification.NotificationModule)(nil)
var _ modules.IHttpModule = (*onlineorder.OnlineOrderModule)(nil)
var _ modules.IHttpModule = (*promotion.PromotionModule)(nil)
var _ modules.IHttpModule = (*purchase.PurchaseModule)(nil)
var _ modules.IHttpModule = (*queue.QueueModule)(nil)
var _ modules.IHttpModule = (*receipt.ReceiptModule)(nil)
var _ modules.IHttpModule = (*report.ReportModule)(nil)
var _ modules.IHttpModule = (*reservations.ReservationModule)(nil)
var _ modules.IHttpModule = (*scheduling.SchedulingModule)(nil)
var _ modules.IHttpModule = (*splitbill.SplitBillModule)(nil)
var _ modules.IHttpModule = (*supplier.SupplierModule)(nil)
var _ modules.IHttpModule = (*table.TableModule)(nil)
var _ modules.IHttpModule = (*tableside.TablesideModule)(nil)
var _ modules.IHttpModule = (*timeclock.TimeClockModule)(nil)
var _ modules.IHttpModule = (*tips.TipsModule)(nil)
var _ modules.IHttpModule = (*training.TrainingModule)(nil)
var _ modules.IHttpModule = (*waste.WasteModule)(nil)
