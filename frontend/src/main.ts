import { createApp } from 'vue'
import App from './App.vue'

import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import 'primeicons/primeicons.css'
import { definePreset } from '@primeuix/themes'

import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import { createI18n } from 'vue-i18n'
import { dt } from '@primevue/themes'
import '@fortawesome/fontawesome-free/css/fontawesome.css'
import '@fortawesome/fontawesome-free/css/regular.css'
import '@fortawesome/fontawesome-free/css/solid.css'
import '@fortawesome/fontawesome-free/css/brands.css'

import { createWebHistory, createRouter } from 'vue-router'

import { createPinia } from 'pinia'
import auth from '@/services/auth'
import Tooltip from 'primevue/tooltip'
import StyleClass from 'primevue/styleclass'
import Ripple from 'primevue/ripple'
import router from '@/router'

declare module 'vue' {
  interface ComponentCustomProperties {
    $auth: typeof auth
  }
}

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      cashier: 'Cashier',
      kitchen: 'Kitchen',
      admin: 'Admin',
      inventory: 'Inventory',
      product: 'Product | Products',
      order: 'Order | Orders',
      order_items: 'Order Items',
      total: 'Total',
      subtotal: 'Subtotal',
      discount: 'Discount',
      egp: 'EGP',
      search: 'Search',
      signout: 'Signout',
      notifications: 'Notifications',
      clear_all: 'Clear All',
      stashed_orders: 'Stashed Orders',
      chats: 'Chats',
      messages: 'Messages',
      write_message: 'Write Message',
      paylater_orders: 'Paylater Orders',
      checkout: 'Checkout',
      category: 'Category | Categories',
      add_component: 'Add Component',
      name: 'Name',
      quantity: 'Quantity',
      unit: 'Unit',
      status: 'Status',
      actions: 'Actions',
      history: 'History',
      list: 'List',
      report: 'Report | Reports',
      settings: 'Settings',
      language: 'Language | Languages',
      sales: 'Sales',
      validation_required: 'This field is required',
      no_results: 'No results found',
      new_orders_appear: 'New orders will appear here automatically',
      cancel: 'Cancel',
      yes: 'Yes',
      delete: 'Delete',
      save: 'Save',
      retry: 'Retry',
      apply: 'Apply',
      close: 'Close',
      back: 'Back',
      next: 'Next',
      error_occurred: 'An error occurred',
      request_failed: 'Failed to perform the request',
      failed: 'Failed',
      done: 'Done!',
      confirmation: 'Confirmation',
      success: 'Success',
      info: 'Info',
      warn: 'Warning',
      danger: 'Danger',
      id: 'Id',
      key: 'Key',
      value: 'Value',
      email: 'Email',
      phone: 'Phone',
      address: 'Address',
      comment: 'Comment',
      username: 'Username',
      password: 'Password',
      confirm_password: 'Confirm Password',
      customer: 'Customer',
      about_nutrix:
        'Nutrix is an open-source restaurant management system that provides a comprehensive solution for managing your restaurant operations.',
      version_info: 'version / commit hash :',
      paid: 'PAID',
      unpaid: 'UNPAID',
      pending: 'PENDING',
      cancelled: 'CANCELLED',
      in_progress: 'INPROGRESS',
      finished: 'FINISHED',
      stashed: 'STASHED',
      refunded: 'REFUNDED',
      normal: 'NORMAL',
      csv: 'csv',
      has_refunds: 'has refunds',
      profile: 'Profile',
      access_denied: 'Access denied',
      relogin: 'Relogin',
      welcome_back: 'Welcome Back',
      sign_in_subtitle: 'Sign in to continue to NutrixPOS',
      enter_username: 'Enter your username',
      enter_password: 'Enter your password',
      sign_in: 'Sign In',
      or: 'or',
      no_account: "Don't have an account?",
      contact_admin: 'Contact admin',
      username_required: 'Username is required',
      password_required: 'Password is required',
      invalid_credentials: 'Invalid username or password',
      login_failed: 'Login failed. Please try again.',
      welcome_to_nutrix: 'Welcome to NutrixPOS',
      no_db_configured: 'No database connection has been configured yet.',
      provide_db_details: 'Please provide your MongoDB connection details below.',
      host: 'Host',
      host_placeholder: 'e.g. localhost',
      port: 'Port',
      port_placeholder: 'e.g. 27017',
      database_name: 'Database name',
      database_placeholder: 'e.g. nutrix',
      username_placeholder: 'e.g. webadmin',
      password_placeholder: 'Your database user password',
      confirm_password_placeholder: 'Enter your password again',
      config_saved_restart:
        'Configuration saved! The server is setup — please restart the application to apply the new settings.',
      connection_successful: 'Connection successful! You can now save the configuration.',
      test_connection: 'Test Connection',
      save_and_continue: 'Save & Continue',
      lets_go: "Let's go",
      host_required: 'Host is required',
      port_invalid: 'A valid port number is required',
      database_required: 'Database name is required',
      passwords_no_match: 'Passwords do not match',
      db_connection_failed: 'Failed to connect to database. Check the credentials.',
      config_save_failed: 'Failed to save configuration. Check the backend logs.',
      create_admin_user: 'Create Admin User',
      no_admin_exists: 'No admin user exists yet.',
      create_admin_subtitle: 'Create your administrator account to get started.',
      admin_placeholder: 'e.g. admin',
      email_placeholder: 'e.g. admin@example.com',
      admin_created_success: 'Admin created successfully! Redirecting to home...',
      create_admin: 'Create Admin',
      email_required: 'Email is required',
      email_invalid: 'Invalid email format',
      password_min_length: 'Password must be at least 6 characters',
      admin_create_failed: 'Failed to create admin user',
      delivery_info: 'Delivery info',
      delivery: 'Delivery',
      anonymous: 'Anonymous',
      drafts: 'Drafts',
      main_details: 'Main details',
      insufficient_availability: 'Insufficient availability',
      only_x_left: 'has only {count} left',
      failed_remove_stashed_order: 'Failed to remove stashed order from db',
      failed_load_stashed_orders: 'Failed to get stashed orders',
      order_stashed_success: 'Order stashed successfully!',
      order_stashed_detail: 'Successfully stashed order!',
      error_stashing_item: 'Error Stashing Item',
      order_finished: 'Order Finished',
      order_ready_to_serve: 'order is ready to be served!',
      order_in_progress: 'Order in progress!',
      drag_drop_upload: 'Drag and drop files here to upload.',
      enable_fixed_cost: 'Enable Fixed Cost',
      fixed_cost: 'Fixed Cost',
      ready_invalid: 'Ready is not a valid number.',
      price_invalid: 'Price is not a valid number.',
      product_deleted_success: 'Product deleted successfully',
      product_delete_failed: 'Failed to delete product',
      confirm_delete_product: 'Are you sure you want to delete this product?',
      product_updated: 'Product Updated',
      product_added: 'Product Added',
      products_load_failed: 'Failed to load products',
      confirm_delete_material: 'Are you sure you want to delete this material?',
      material_deleted_success: 'Material deleted successfully',
      material_delete_failed: 'Failed to delete material',
      material_edited_success: 'Successfully edited material data',
      material_settings_saved: 'Material settings saved',
      confirm_delete_entry: 'Are you sure you want to delete this entry?',
      entry_deleted: 'Entry deleted!',
      entry_saved: 'Entry saved successfully!',
      component_saved: 'Component saved successfully!',
      category_added_success: 'Category added successfully',
      category_add_failed: 'Failed to add category',
      category_deleted_success: 'Category deleted successfully',
      category_delete_failed: 'Failed to delete category',
      category_updated_success: 'Category updated successfully',
      category_update_failed: 'Failed to update category',
      confirm_delete_category: 'Are you sure you want to delete this category?',
      categories_load_failed: 'Failed to fetch categories',
      order_cancelled_success: 'Order cancelled successfully',
      order_cancel_failed: 'Failed to cancel order',
      confirm_cancel_order: 'Are you sure you want to cancel this order?',
      users_load_failed: 'Failed to load users',
      user_added: 'User Added',
      user_deleted_success: 'User deleted successfully',
      user_delete_failed: 'Failed to delete user',
      confirm_delete_user: 'Are you sure you want to delete this user?',
      enter_new_password: 'Please enter a new password',
      password_changed_success: 'Password changed successfully',
      password_change_failed: 'Failed to change password',
      customer_deleted_success: 'Customer deleted successfully',
      customer_delete_failed: 'Failed to delete customer',
      confirm_delete_customer: 'Are you sure you want to delete this customer?',
      customer_updated: 'Customer Updated',
      customer_added: 'Customer Added',
      customers_load_failed: 'Failed to load customers',
      shop_mode: 'Shop Mode',
      shop_mode_description:
        'Controls which features are available. Restart the app after changing.',
      retail: 'Retail',
      settings_updated: 'Settings updated successfully!',
      settings_update_failed: 'Error updating settings!',
      settings_load_failed: 'Failed to load settings. Please try again.',
      language_load_error: 'Error loading language',
      languages_load_failed: 'Failed to load languages',
      hubsync: 'Hubsync',
      server_host: 'Server Host',
      token: 'Token',
      sync_interval: 'Sync Interval (seconds)',
      buffer_size: 'Buffer Size',
      data_saved_success: 'Data saved successfully',
      select_language: 'Select a Language',
      fill_all_fields: 'Please fill all fields',
      add_with_comment: 'Add with comment',
      product_details: 'Product details',
      confirm_finish_order: 'Are you sure you want to finish the order?',
      order_finished_detail: 'Good job!',
      finish_order_failed: 'Failed finishing order!',
      custom_data: 'Custom data',
      edit_custom_data: 'Edit custom data',
      add_custom_data: 'Add custom data',
      tip_added: 'Tip added',
      tip_removed: 'Tip removed',
      confirm_delete_custom_data: 'Are you sure you want to delete this custom data entry?',
      custom_data_updated: 'Custom data updated',
      money_collected: 'Money collected',
      print_signal_sent: 'Print signal sent',
      items: 'Item(s)',
      return: 'return',
      returning_item: 'Returning',
      money_to_refund: 'Money To refund',
      destination: 'Destination:',
      dispose: 'Dispose',
      waste: 'Waste',
      custom_materials_products: 'Custom (materials / products)',
      materials: 'Materials',
      item_refunded: 'Item refunded',
      refund_error: 'Refund failed',
      order_paid_success: 'Successfully paid order',
      about_nutrix_short: 'Nutrix is an open-source restaurant management system',
    },
  },
})

const Noir = definePreset(Aura, {
  components: {
    progressspinner: {
      colorScheme: {
        light: {
          root: {
            colorOne: '{primary.900}',
            colorTwo: '{primary.900}',
            colorThree: '{primary.900}',
            colorFour: '{primary.900}',
          },
        },
        dark: {
          root: {
            colorOne: '{primary.900}',
            colorTwo: '{primary.900}',
            colorThree: '{primary.900}',
            colorFour: '{primary.900}',
          },
        },
      },
    },
  },
  semantic: {
    primary: {
      50: '{zinc.50}',
      100: '{zinc.100}',
      200: '{zinc.200}',
      300: '{zinc.300}',
      400: '{zinc.400}',
      500: '{zinc.500}',
      600: '{zinc.600}',
      700: '{zinc.700}',
      800: '{zinc.800}',
      900: '{zinc.900}',
      950: '{zinc.950}',
    },
    colorScheme: {
      light: {
        primary: {
          color: '#14977B',
          inverseColor: '#a5c22f',
          hoverColor: '#14977B',
          activeColor: '#14977B',
        },
        highlight: {
          background: '#DEDB69',
          focusBackground: '#DEDB69',
          color: '#173350',
          focusColor: '#173350',
        },
      },
      dark: {
        primary: {
          color: '#a5c22f',
          inverseColor: '#14977B',
          hoverColor: '#a5c22f',
          activeColor: '#a5c22f',
        },
        highlight: {
          background: '#fff6c7',
          focusBackground: '#FFDC00',
          color: '#173350',
          focusColor: '#173350',
        },
      },
    },
  },
})

const app = createApp(App).use(createPinia())
app.config.globalProperties.$auth = auth

app
  .use(router)
  .use(PrimeVue, {
    theme: {
      preset: Noir,
      options: {
        prefix: 'p',
        darkModeSelector: '.my-app-dark',
        cssLayer: false,
      },
    },
  })
  .use(ToastService)
  .use(ConfirmationService)
  .use(i18n)
  .directive('tooltip', Tooltip)
  .directive('styleclass', StyleClass)
  .directive('ripple', Ripple)
  .mount('#app')
