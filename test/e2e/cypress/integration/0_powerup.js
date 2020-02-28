const url = Cypress.env('PROXEUS_PLATFORM_DOMAIN') ? Cypress.env(
  'PROXEUS_PLATFORM_DOMAIN') : 'http://localhost:1323'

describe(`Power up on ${url}`, () => {

  before(() => {
    cy.visit(`${url}`)
  })

  it('should have a powerup button', () => {
    const button = cy.get('button.btn-primary').eq(2)
    button.should($button => {
      expect($button).to.be.visible
      expect($button).to.contain.text('Power up')
    })
  })

  it('should bring to home page after saving with default values', () => {
    cy.get('button.btn-primary').eq(2).click()
  })

})
