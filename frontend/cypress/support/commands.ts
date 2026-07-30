/// <reference types="cypress" />

Cypress.Commands.add('login', (username: string, password: string) => {
  cy.visit('/login')
  cy.get('input').first().type(username)
  cy.get('input').last().type(password)
  cy.get('button[type="submit"], button').contains(/login|prijava/i).click()
})

/* eslint-disable @typescript-eslint/no-namespace */
declare global {
  namespace Cypress {
    interface Chainable {
      login(username: string, password: string): Chainable<void>
    }
  }
}
/* eslint-enable @typescript-eslint/no-namespace */
export {}
