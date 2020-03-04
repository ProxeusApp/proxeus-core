const url = Cypress.env('PROXEUS_PLATFORM_DOMAIN') ? Cypress.env(
  'PROXEUS_PLATFORM_DOMAIN') : 'http://localhost:1323'

const password = 'test-password'

describe(`User signup & login at ${url}`, () => {

  describe('User creation', () => {
    before(() => {
      cy.visit(`${url}`)
    })

    const adminEmailAddress = `test-admin-email${Math.floor(
      (Math.random() * 10000) + 1)}@proxeus.com`
    const emailAddress = `test-user-email${Math.floor(
      (Math.random() * 10000) + 1)}@proxeus.com`

    let signupResponse // Used to populate with token once registered, in order to activate the user

    context('page layout', () => {

      it('should have a signup button on the top', () => {
        const button = cy.get('.btn-primary#signup')
        button.should($button => {
          expect($button).to.be.visible
          expect($button).to.have.text('Sign up')
        })
        cy.get('.btn-primary').first().click()
      })

      it('should have another signup button', () => {
        const button = cy.get('#signupcontent')
        button.should($button => {
          expect($button).to.be.visible
          expect($button).to.contain.text('Sign up')
        })
      })

      it('should go to signup page by clicking on first signup button', () => {
        cy.get('.btn-primary#signup').click()
        cy.url().should('include', '/register')
      })

      it('should go to signup page by clicking on second signup button', () => {
        cy.get('#signupcontent').click()
        cy.url().should('include', '/register')
      })
    })

    context('admin signup page', () => {
      before(() => {
        cy.visit({ url: `${url}/register`, failOnStatusCode: false })
      })

      it('should have a mail input', () => {
        cy.get('#inputEmail')
      })

      it('should sign up admin', () => {
        cy.wait(2000)
        cy.server()
        cy.route({
            method: 'POST',
            url: '/api/register',
            onResponse: (xhr) => {
              signupResponse = xhr.response.headers['x-test-token'];
              console.log(xhr.response);
            },
          },
        ).as('sign-admin')

        cy.get('#inputEmail').type(`${adminEmailAddress}{enter}`)
        cy.get('#frontend-app').should('contain', 'Email sent')
      })

      it('should activate the admin via valid token', () => {
        cy.request({
          method: 'POST',
          url: `${url}/api/register/${signupResponse}`,
          body: {
            password: password,
          },
        })
      })
    })

    context('user signup page', () => {
      before(() => {
        cy.visit({ url: `${url}/register`, failOnStatusCode: false })
      })

      it('should sign up user', () => {
        cy.server()
        cy.route({
            method: 'POST',
            url: '/api/register',
            onResponse: (xhr) => {
              signupResponse = xhr.response.headers['x-test-token']
            },
          },
        )

        cy.get('#inputEmail').type(`${emailAddress}{enter}`)
        cy.get('#frontend-app').should('contain', 'Email sent')
      })

      it('should NOT activate the user via a wrong token', async () => {
        const response = await cy.request({
          method: 'POST',
          url: `${url}/api/register/00000000-0000-0000-0000-00099bd95ac7`,
          body: {
            password: password,
          },
          failOnStatusCode: false,
          timeout: 2000,
        })

        expect(response.status).not.to.eq(200)
      })

      it('should activate the user via valid token', async () => {
        await cy.request({
          method: 'POST',
          url: `${url}/api/register/${signupResponse}`,
          body: {
            password: password,
          },
        })
      })
    })

    context('login', () => {
      before(() => {
        cy.visit(`${url}/login`)
      })

      it('should fail with wrong credentials', () => {
        login(emailAddress, 'wrong-password')

        cy.get('.text-danger').
          should('contain',
            'You have entered an invalid username or password')
      })

      it('should login successfully', () => {
        login(emailAddress, password)

        cy.get('.navbar-text').should('contain', 'Workflows')
      })

    })

    context('profile', () => {

      before(() => {
        cy.visit(`${url}/admin/workflow`)
        login(emailAddress, password)
      })

      it('should open the account box menu', () => {
        cy.get('i.trp-photo.material-icons').click()

        cy.get('.btn-primary').should('contain', 'Account')
      })

      it('should open the account form', () => {
        cy.get('.btn-primary').first().click()
      })

      it('should have a DELETE button', () => {
        cy.get('.btn-danger')
      })

      it('should ask again before removing account', () => {
        cy.get('.btn-danger').click()
      })

      it('should remove account and logout when clicking again', () => {
        cy.get('.btn-danger').click()
        cy.url().should('eq', `${url}/login`)
      })
    })
  })

})

function login (email, password) {
  cy.get('#inputEmail').clear().type(email)
  cy.get('#inputPassword').clear().type(password)

  cy.get('.btn-primary').eq(0).click()
}
