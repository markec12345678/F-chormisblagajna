import { createRouter, createWebHistory } from 'vue-router'
import auth from '@/services/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/setup',
      component: () => import('@/pages/Setup.vue'),
    },
    {
      path: '/admin-setup',
      component: () => import('@/pages/AdminSetup.vue'),
    },
    {
      path: '/login',
      component: () => import('@/pages/Login.vue'),
    },
    {
      path: '/no-access',
      component: () => {
        return import('@/pages/NoAccessView.vue')
      },
    },
    {
      path: '/',
      alias: ['/home'],
      component: () => {
        if (!auth.isAuthenticated.value) {
          window.location.href = '/login'
          return import('@/pages/Login.vue')
        }
        if (auth.hasRole('admin') || auth.hasRole('superuser') || auth.hasRole('cashier')) {
          return import('@/pages/Home.vue')
        }
        return import('@/pages/NoAccessView.vue')
      },
    },
    {
      path: '/kitchen',
      component: () => {
        if (!auth.isAuthenticated.value) {
          window.location.href = '/login'
          return import('@/pages/Login.vue')
        }
        if (auth.hasRole('admin') || auth.hasRole('superuser') || auth.hasRole('chef')) {
          return import('@/pages/Kitchen.vue')
        }
        return import('@/pages/NoAccessView.vue')
      },
    },
    {
      path: '/admin',
      component: () => {
        if (!auth.isAuthenticated.value) {
          window.location.href = '/login'
          return import('@/pages/Login.vue')
        }
        if (auth.hasRole('admin') || auth.hasRole('superuser')) {
          return import('@/pages/Admin.vue')
        }
        return import('@/pages/NoAccessView.vue')
      },
      children: [
        {
          path: '',
          redirect: { path: '/admin/dashboard' },
        },
        { path: 'dashboard', component: () => import('@/pages/Dashboard.vue') },
        { path: 'inventory', component: () => import('@/pages/Inventory.vue') },
        { path: 'sales', component: () => import('@/pages/Sales.vue') },
        { path: 'products', component: () => import('@/pages/Products.vue') },
        { path: 'categories', component: () => import('@/pages/Categories.vue') },
        { path: 'orders', children: [{ path: '', component: () => import('@/pages/Orders.vue') }] },
        { path: 'settings', component: () => import('@/pages/Settings.vue') },
        { path: 'customers', component: () => import('@/pages/Customers.vue') },
        {
          path: 'users',
          component: () => {
            if (auth.hasRole('superuser')) {
              return import('@/pages/Users.vue')
            }
            return import('@/pages/NoAccessView.vue')
          },
        },
        {
          path: 'hubsync',
          component: () => {
            if (auth.hasRole('superuser')) {
              return import('@/pages/Hubsync.vue')
            }
            return import('@/pages/NoAccessView.vue')
          },
        },
        { path: 'tables', component: () => import('@/pages/Tables.vue') },
        { path: 'branches', component: () => import('@/pages/Branches.vue') },
        { path: 'split-bills', component: () => import('@/pages/SplitBill.vue') },
        { path: 'scheduling', component: () => import('@/pages/Scheduling.vue') },
        { path: 'reservations', component: () => import('@/pages/Reservations.vue') },
        { path: 'promotions', component: () => import('@/pages/Promotions.vue') },
        { path: 'loyalty', component: () => import('@/pages/Loyalty.vue') },
        { path: 'reports', component: () => import('@/pages/Reports.vue') },
        { path: 'locations', component: () => import('@/pages/MultiLocationDashboard.vue') },
        { path: 'giftcards', component: () => import('@/pages/GiftCards.vue') },
        { path: 'transfers', component: () => import('@/pages/InventoryTransfers.vue') },
        { path: 'fiscal', component: () => import('@/pages/FiscalDashboard.vue') },
        { path: 'floorplan', component: () => import('@/pages/FloorPlan.vue') },
        { path: 'tips', component: () => import('@/pages/EmployeeTips.vue') },
        { path: 'accounting', component: () => import('@/pages/Accounting.vue') },
        { path: 'menu-engineering', component: () => import('@/pages/MenuEngineering.vue') },
        {
          path: 'employee-performance',
          component: () => import('@/pages/EmployeePerformance.vue'),
        },
        { path: 'suppliers', component: () => import('@/pages/Suppliers.vue') },
        { path: 'expenses', component: () => import('@/pages/Expenses.vue') },
        { path: 'timeclock', component: () => import('@/pages/TimeClock.vue') },
        { path: 'waste', component: () => import('@/pages/WasteTracking.vue') },
        { path: 'online-orders', component: () => import('@/pages/OnlineOrders.vue') },
        { path: 'chat', component: () => import('@/pages/StaffChat.vue') },
        { path: 'audit-log', component: () => import('@/pages/AuditLog.vue') },
        { path: 'receipt', component: () => import('@/pages/ReceiptCustomization.vue') },
        { path: 'inventory-alerts', component: () => import('@/pages/InventoryAlerts.vue') },
        { path: 'multi-payment', component: () => import('@/pages/MultiPayment.vue') },
        { path: 'customer-display', component: () => import('@/pages/CustomerDisplay.vue') },
        { path: 'tableside', component: () => import('@/pages/TablesideOrdering.vue') },
        { path: 'staff-training', component: () => import('@/pages/StaffTraining.vue') },
        { path: 'customer-feedback', component: () => import('@/pages/CustomerFeedback.vue') },
        { path: 'delivery', component: () => import('@/pages/DeliveryManagement.vue') },
        { path: 'kiosk', component: () => import('@/pages/KioskSettings.vue') },
        { path: 'languages', component: () => import('@/pages/Languages.vue') },
        { path: 'marketing', component: () => import('@/pages/MarketingCampaigns.vue') },
        { path: 'purchase-orders', component: () => import('@/pages/PurchaseOrders.vue') },
        { path: 'queue', component: () => import('@/pages/QueueWaitlist.vue') },
        { path: 'ai', component: () => import('@/pages/AI.vue') },
      ],
    },
    {
      path: '/order',
      component: () => import('@/pages/OnlineOrderPortal.vue'),
    },
    {
      path: '/display/:id',
      component: () => import('@/pages/DisplayView.vue'),
    },
    {
      path: '/tableside/menu/:token',
      component: () => import('@/pages/TablesideMenu.vue'),
    },
    {
      path: '/profile',
      component: () => {
        if (!auth.isAuthenticated.value) {
          window.location.href = '/login'
          return import('@/pages/Login.vue')
        }
        return import('@/pages/Profile.vue')
      },
    },
  ],
})

export default router
