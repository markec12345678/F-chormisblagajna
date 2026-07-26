describe('POS Application', () => {
  describe('Login', () => {
    it('shows login page', () => {
      cy.visit('/login')
      cy.get('input[type="text"], input[placeholder*="user"], input[placeholder*="User"]').should('exist')
    })
  })

  describe('Navigation', () => {
    it('has admin sidebar links', () => {
      cy.visit('/admin')
      cy.get('aside, nav, .sidebar, [class*="sidebar"]').should('exist')
    })
  })

  describe('Tables page', () => {
    it('loads tables page', () => {
      cy.visit('/admin/tables')
      cy.get('.datatable-stub, table, .p-datatable').should('exist')
    })
  })

  describe('AI Search page', () => {
    it('loads AI search page', () => {
      cy.visit('/admin/ai')
      cy.get('input, .p-inputtext').should('exist')
    })
  })

  describe('Categories', () => {
    it('loads categories page', () => {
      cy.visit('/admin/categories')
      cy.get('table, .p-datatable, .datatable-stub').should('exist')
    })
  })
})
