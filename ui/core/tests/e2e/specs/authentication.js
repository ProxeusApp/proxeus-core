
describe('Authentication', () => {
  it('Logs in and out of Proxeus Core', () => {
    cy.visit('/proxeuslogin')
    cy.get('input[name="email"]').type('admin')
    cy.get('input[name="password"]').type('supernimda{enter}')
    cy.contains('h1', 'Workflows')
    cy.contains('Account').click()
    cy.contains('Logout').click()
    // Back on frontpage
    cy.contains('Document Verification')
  })
})
