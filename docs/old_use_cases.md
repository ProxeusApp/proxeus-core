# Use Cases

Proxeus plays the role of a bridge connecting two worlds, allowing you to digitize traditional processes and register information on a blockchain, as well as to make blockchain activity visible by generating human-readable documents.

Legal
=====

**Incorporation of Swiss companies**

![](_media/old_proxeus/legal/1.svg)

## Background

The process of legally registering a new company in Switzerland heavily relies on printed documents that have to be exchanged between several parties, and can take up to six weeks. Correctly filling out documents and sending right ones to the right parties in the right order alone can cost weeks to people unfamiliar with the process.

The gathering and exchange of information between the participants is very time consuming and error-prone. The main document can be 30 pages long and every detail needs to be thoroughly verified by a notary before the documentation is handed over to the official business registry. As there are

## Use case exploration

As part of a challenge by [digitalswitzerland](https://digitalswitzerland.com/), Proxeus joined a project initiated in 2017 by IBM, E&Y and Swisscom, aiming to change the status quo and accelerate the registration of Swiss companies. The project team consisted of lawyers, notaries, banks, authorities (register of commerce), IT service providers and software developers.

Following several workshops where the various processes were analyzed, discussed and streamlined. A prototype was built and several improvement loops were done until the point when the group felt confident enough to present the solution publicly. In May 2018, a public demonstration showed that the total time to incorporate a company can be reduced to less than two hours by using the Proxeus document workflow automation engine and a blockchain to coordinate processes more effectively.

## Implementation

The solution consists of a Proxeus platform with workflows for the different legal forms that can be registered, starting with simple forms and adding others later. The workflows comprises a series of smart forms that validate all inputs immediately. This highly simplifies the processes for data entry, document creation and data exchange between all the parties that are involved. Information that had to be entered redundantly up to 30 times is now only entered once in a Proxeus form. Input validation that had to be done manually before, happens automatically now. Tasks that previously had to be tackled sequentially, due to the need to move paper dossiers around, can now be completed in parallel.

![](_media/old_proxeus/legal/2.png)

Excerpt from one of the user forms of the incorporation process

The workflow also includes a conditional alternative path, depending on if the founder has already opened a capital deposit account at a bank. If he has not, two additional forms are shown and one more document is produced.

![](_media/old_proxeus/legal/4.png)

Conditional alternative path


At the end of the workflow, all needed documents are created and their hash is registered on the Ethereum blockchain, making the document tamper-proof and verifiable. While Proxeus provides a verification tool for its users’ convenience, the entries are made on the public Ethereum blockchain and may also be checked using block explorers like Etherscan or custom tools.

![](_media/old_proxeus/legal/6.png)

One of the workflow versions used to incorporate a company

At the end of the workflow, a custom node built for Proxeus’ integration layer interacts with a service developed by the project partner IBM. Through this service, Proxeus is connected to a votation contract on Hyperledger Fabric, allowing the relevant stakeholders (entrepreneur, lawyer, bank, notary, public register) to verify and sign the documents produced. The bank states that the capital money has indeed been paid; the notary confirms that the necessary documents have been provided, read over, and approved; and the commercial register performs a final check that everything is accurate and lawful. After successful collection of all confirmations, the filing is officially registered with the Commercial Register and Official Gazette of Commerce.

## Result

The Proxeus platform set up specially for the business registry project was well-received by all parties. Using the workflows proved intuitive for users of all kinds of professional backgrounds. When the first real business incorporation was done using the platform, it took less than two hours, crushing the initial objective of 48 hours.

<strong> This achievement was a Swiss – and likely world’s – first! </strong>

The prototype worked so well that the test phase has then been prolonged to register up to 100 companies with Proxeus. By now, dozens of companies have been registered through Proxeus workflows and the process has been fine-tuned several times by the project team. In 2020 the platform is still live and being used and maintained by the project participants. The next step for the project, which has already began, is the onboarding of further notaries and lawyers, who will at their turn incorporate companies through Proxeus.


## Feedback

In an [interview](https://medium.com/proxeus/proxeus-helps-speeding-up-swiss-business-incorporations-dd0eed421576) in August 2019, the lawyer Philippe Kaiser from Kaiser Odermatt & Partner, who was involved in the project, commented on the Proxeus solution:

>“The tool highly simplifies the processes for data entry, document creation and data exchange between all the parties that are involved. [...] We can certainly imagine many other use cases. The Proxeus tool is very simple and intuitive to use. In particular, the intuitive handling of the tool opens up many possibilities and easy access for users that aren’t IT experts.”

The project received massive media attention and made it into several big newspapers, magazines and TV shows. It sparked a desire for innovation in several of the cantonal business registers of Switzerland.


## Insights

The project was initially intended as a prototype but was consolidated into a more stable solution as the project participants decided to bring it to MVP level, generating a feedback stream that was useful for improving several aspects of the product. The feedback collected allowed to design and scope a series of sharing features, and led to the improvement of the import/export capabilities of Proxeus, going beyond the requirements set in the original project whitepaper.

The project also led to the realization that a higher degree of automation would require more integration, and that because every participating party has existing IT systems in place, a well-designed integration layer should be made available in order to allow for Proxeus platforms to scale in the future.


## Limitations

The project was initially intended as a prototype but was consolidated into a more stable solution as the project participants decided to bring it to MVP level, generating a feedback stream that was useful for improving several aspects of the product. The feedback collected allowed to design and scope a series of sharing features, and led to the improvement of the import/export capabilities of Proxeus, going beyond the requirements set in the original project whitepaper. 


## Try it yourself

If you enjoyed reading the documentation of this project and would like to try building something similar, we suggest following the steps below. For the most part you’ll only need decent skills in using Proxeus, but for the full scope some programming knowledge is required.

Here is how you can create a workflow - using only Proxeus and no programming at all.

  - Understand the requirements. What workflow output do you expect? What documents should be registered? What role should signatures play?

  - Requirements analysis and specification:

    - Which part of which business process would you like to digitize with your Proxeus prototype?

    - Who are the stakeholders? What are the different tasks? What are the roles?

    - Which documents shall be produced and what information is needed to create them?

    - On which blockchain do you want to register the resulting documents? Do you require a smart contract or would simple transactions suffice?

-   Set up Proxeus. You can run your own instance of Proxeus on a server or locally on your computer. The complete guide to setting up your own instance is available [here](https://doc.proxeus.org/#/README). It is recommended that you deploy your own smart contract following our instructions in the guide and using the template in our [GitHub](https://github.com/ProxeusApp/proxeus-contract.git).

-   Create Proxeus workflows. When you scoped your project in step 1, you figured out which documents need to be produced. After you’ve set up your own instance of Proxeus, you can now configure a workflow for each document. Workflows comprise of data entry forms and document templates. Some users prefer to design the document templates first and then to derive the necessary data inputs in the forms from there. The other way around is also fine.
-   Prepare the document template(s). What should the document design look like? What should placeholders be used for? How should the information be formatted (e.g. what sections should be shown vs. hidden in the final output)?
Excerpt from one of the templates of the incorporation process

![](_media/old_proxeus/legal/7.png)

<em>Excerpt from one of the templates of the incorporation process</em>

-   Prepare the user forms. Adapt the requirements to the audience (e.g. athletes cannot be expected to be very tech-savvy) - entering the data has to be intuitive and instill confidence. All Proxeus form elements support help texts - they even accept HTML and links to further information - use them to clarify every step.

-   Optimize your workflow for efficient, intuitive data entry. Help the user avoid mistakes by explaining each field with help texts and by validating the fields. Proxeus can support you by checking if the input has the expected format, for example a valid email address or date. Create smart fields that only appear when certain conditions are met. For example, you can use smart placeholders like this:
‍
```
{% if (input.number > 0) %} Number is higher than 0 {% else %} 0 {% endif %}‍
```


- Interact with blockchains or other systems. Thanks to the integration interfaces of Proxeus, you can simply add custom code into your workflows. Read data from an ERP, send an email confirmation, create a transaction in blockchains other than Proxeus’ native platform Ethereum - all this can be done through custom nodes. The technical documentation of our external nodes library can be found on our GitHub.
  

To help you kickstart your ideas, we’ve provided several examples on this website. In the project for business incorporations we made Proxeus communicate with an API for IBM’s Hyperledger Fabric framework. A votation contract was used to collect the signatures of the participants, each of them confirming that their tasks have been completed. This specific piece of code cannot be shared, but can be recreated using Hyperledger’s documentation - or connecting Proxeus to a voting contract on any blockchain of your choice.

- Set signatures as required - once given by users they are publicly visible on the blockchain and verifiable by anyone.

-  Configure the workflow by connecting the smart template(s) and the forms in a workflow. Set a XES price to the workflow (if applicable) and share the form with the platform users.

-  Test and improve your workflow. Nothing is perfect from the start. Even when your workflow is already in production and in use, you can simply clone it and release an improved version to your users. If the changes are compatible, you can also just upgrade the existing workflow and all running instances will use the new definition automatically.

-  Create an instruction page to onboard users. Our handbook and tutorial videos are also here to help.

-  Integrate a verification tool into your organization’s website to simplify the verification process. The verification tool of your Proxeus platform can be found directly via the /validation path; for example https://morrison.proxeus.org. It can also be embedded using a standard iFrame HTML code.

<iframe width="100%" height="650" src="https://morrison.proxeus.org/" frameborder="0" marginwidth="0" marginheight="0" scrolling="yes"></iframe>

Tamper-proof
============

**University Diplomas**

![](_media/old_proxeus/education/1.svg)

## Background

Suspicion of fraud leads an increasing number of companies to set up verification procedures, creating a backlog of work for education institutions. Blockchain technology makes it possible to issue trusted digital documents that can be verified by employers independently.

## Use case exploration

Over the course of Summer 2018, a discussion was engaged with Professor Dr. Schär, Managing Director of the Center for Innovative Finance at the University of Basel in Switzerland, who presented a basic solution design showing how university diplomas could be certified on the Ethereum blockchain. A project was initiated to prototype his idea within a one-week timeframe using the Proxeus framework.

## Implementation

The project comprised the deployment of a new Proxeus instance, the collection of requirements, the design of a diploma template, the configuration of the diploma creation workflow as well as some fine-tuning and front-end improvements. The team started by creating the template, which allowed to review and validate the required data entry fields based on the university feedback. The diploma looked similar to the following example:

![](_media/old_proxeus/education/2.png)

The next step was to create a series of input forms allowing the university staff to enter the relevant student’s information (personal identification data and grade).

![](_media/old_proxeus/education/3.png)

  

Once the templates and forms were done, bringing everything in the right order was just a quick drag & drop exercise. The template to be filled comes first, then the data entry forms that guide the user through the process.

![](_media/old_proxeus/education/4.png)

A workflow with two forms to fill one template

Finally, a verification interface was built into the Core allowing non-technical users of Proxeus to easily verify the authenticity of documents created using the platform.

![](_media/old_proxeus/education/5.png)

## Result

Following a successful first test, Prof. Dr. Schär proposed to go beyond the initial project scope and issue real certificates for his students on the Ethereum mainnet. The university staff was provided with a ready-to-use workflow allowing them to issue course certificates and a few weeks later, all 126 students completing the “Bitcoin, Blockchain and Crypto Assets” course of the University of Basel received certificates registered on the Ethereum mainnet blockchain. This marked the first time that a university in Switzerland secured the authenticity and integrity of an academic certificate on a blockchain, and an important milestone for Proxeus.  
  
The students and everyone with whom they decide to share the certificates can verify the authenticity of the document instantly, 24/7, on a Proxeus verification interface embedded into the university website.  
  
The new verification tool was a direct result of this project and was later used in almost every project. It was also made available for embedding into websites. Check out the step-by-step guide at the end of this article to find out how.  

## Feedback

Prof. Dr. Fabian Schär commented on the project in a university press release:_  
  
“This blockchain solution developed in cooperation with Proxeus greatly improves the process of verifying documents and represents an important step towards forgery-proof academic certificates,” said Prof. Dr. Fabian Schär, Managing Director of the Center for Innovative Finance. “Two tools have been created – one that the university staff can use to create the certificate and register it to the blockchain with its unique hash, and a second, public tool which allows anyone to verify that the electronic document has been issued by the university.”  
  
“Fraudulent documents are a problem in academia just as it is in any field,” he said. “By securing credentials on the Blockchain, we provide an extra layer of security for graduates and potential employers. These credentials can’t be faked, and can be easily verified online. It will introduce a new paradigm of security and offer value to all parties - employers don’t lose time checking credentials, graduates have an edge, and the institutions themselves reduce their reputational risk and a significant administrative burden.”_‍

## Insights

Most university’s requirements received as part of the project could be covered by Proxeus out of the box. Important learnings were made as to how blockchain data can be made accessible to everyone and a visual drag and drop interface was included into the Proxeus Core software. Observations were also made regarding the API integration possibilities with external systems (such as the university ERP) in order to allow for projects to scale. A powerful I/O layer was designed and tested in various further projects in the course of 2019.  
  
Today, the tamper-proof diploma use case can be replicated in under an hour with Proxeus out of the box and it demonstrates what Proxeus intends to facilitate: making the creation of blockchain applications possible in a few clicks and bringing the same level of simplicity as WordPress does for websites.

## Limitations

Proxeus is a workflow engine enabling users to digitize their processes and make document templates available to third-parties decentrally. Its capabilities are therefore focused on the production of new certificates, not on the processing of pre-existing documents, which would require a different solution design. A custom module (out of scope for this project) could be developed to import data directly from a university ERP. One of the key requirements of the university was that the documents would be secured additionally using a so-called "salt". A salt is a unique code increasing the cryptographic difficulty of guessing the contents of a document, as someone in possession of an almost identical document (e.g. only the name and grade change) could theoretically modify those parts and try to reverse-engineer how a hash was made. This feature was too specific to be integrated into the generic Proxeus solution, but can be added manually into the document as invisible (white) text as we have demonstrated.

## Try it out yourself

A very similar version of the workflow developed in this project is available on the [demo platform](https://morrison.proxeus.org/) on Proxeus.com.  

You enjoyed reading the documentation and would like to build a similar workflow? A detailed guide is available [here](https://docs.google.com/document/d/1Gl6R1t0LYRK6kARScx5-vqpiKtd6xqcH2yPPC3wWMr4/preview). All you need is an instance of Proxeus (for example our [demo platform](http://demo.proxeus.com/)) and basic users skills.  

If you are just looking to understand the principles behind the solution, here is a recap of the steps to build it - using only Proxeus and no programming at all.

*   Understanding the requirements. We first discussed what exactly was expected from the workflow, the document certification on the blockchain and the verifiability.
*   Set up Proxeus. You can run your own instance of Proxeus on a server or locally on your computer - or you could use someone else’s instance. The complete guide to setting up your own instance is available [here](http://doc.proxeus.com/#/README). It is recommended that you deploy your own smart contract following our instructions in the guide and using the [template in our GitHub](https://github.com/ProxeusApp/proxeus-contract).
*   Prepare the document template(s). What should the document design look like? What should placeholders be used for? How should the information be formatted? In the use case described above, a user can for example pick a date from the calendar and the template will display the weekday in alphanumeric form (e.g. “Friday”).  
    Here are the direct links to our help materials:  
    \- [General handbook](https://docs.google.com/document/d/1C3B1oNY6lOv8Q_AvbKhwlySrS6qTiRl3raPLV6OXr7w/preview)  
    \- [Template handbook](https://docs.google.com/document/d/1-vJsTrU3w8dEcDr3-nV5owtxqHWSjzEf2uk6m9-cMIs/preview)  
    \- [Education example template](https://docs.google.com/document/d/1Gl6R1t0LYRK6kARScx5-vqpiKtd6xqcH2yPPC3wWMr4/edit)
*   Prepare the user forms. As the university staff will need to use the workflow for each student, data entry has to be efficient and comfortable. We decided to split it into two forms to avoid scrolling: one form for the student’s data and the second for the grade.
*   Configure the workflow by connecting the smart template(s) and the forms in a workflow. Set a XES price to the workflow (if applicable) and share the form with the platform users.
*   Set signatures as required - once given by users they are publicly visible on the blockchain and verifiable by anyone.
*   Test and improve your workflow. Nothing is perfect from the start. Even when your workflow is already in production and in use, you can simply clone it and release an improved version to your users. If the changes are compatible, you can also just upgrade the existing workflow and all running instances will use the new definition automatically.
*   Onboard your users. In the use case described above only a very brief training was required to make sure the university staff would be able to use the platform safely. Our handbook and tutorial videos are always there as an aid in case something is forgotten. Basically, the users only need to know the following steps:  
    \- Log in  
    \- Click “New document” in the “Documents” tab  
    \- Select the workflow you’ve created  
    \- Navigate through the workflow  
    \- Complete the workflow and download the document

![](_media/old_proxeus/education/6.png)

_Starting a new workflow to create a blockchain-secured diploma_

9\. Integrate a verification tool into your organization’s website to simplify the verification process. The verification tool of your Proxeus platform can be found directly via the /validation path; for example [https://morrison.proxeus.org/](https://morrison.proxeus.org/). It can also be embedded using a standard iFrame HTML code.  
  
Your iFrame HTML code to embed it into your website would then look like this:

Sports
======

**Anti-doping process**

![](_media/old_proxeus/sports/1.svg)

## Background


Professional sport is a matter of national pride and results can have considerable financial consequences. For this reason, this field is highly prone to fraud, which has led to the occurrence of several public scandals over the last years. A series of blockchain use cases can be imagined in this area to increase process efficiency and auditability.  
‍  
The [International Testing Agency](http://ita.sport/) (ITA) was officially founded in January 2018 following a proposal by the International Olympic Committee (IOC) to make anti-doping independent from national and international sports organisations. The organisation has been established by the IOC as an independent, not-for-profit foundation under Swiss law and is headquartered in Lausanne, Switzerland.  
  
One of the key roles of anti-doping organizations is to prevent the materialization of a series of important risks:  
  
\- testing plans (who gets tested when by who) can be leaked  
\- test results can be tampered with  
\- storage servers can be compromised  
\- information shared between parties can be intercepted  
  
ITA is dealing against powerful forces and must therefore constantly look to improve itself and address potential vulnerabilities in its processes. In this context, a project has been discussed and setup to explore the possible applications of blockchain in the anti-doping space.

## Use case exploration


If an athlete has medical reasons to take a prohibited substance, he/she can apply to receive a certificate granting them an exemption for therapeutic use (TUE). Such exemption is given after consideration of an independent medical board and for a set prohibited substance, in a certain dosage, within a limited period.  
‍  
The project consisted of prototyping the digitization of the TUE request workflow (currently managed via a .pdf form sent by the athlete via email) using the [Proxeus](http://www.proxeus.com/) framework, thus increasing process efficiency and auditability.

## Implementation


After the initial project scoping, the project consisted of creating and configuring a workflow based on the existing [TUE request template](https://ita.sport/wp-content/uploads/2019/08/TUE-Form.pdf), designing a corresponding template and setting up a signature process.  

The TUE request was implemented as a multi-page .odt smart template with complex variables:

![](_media/old_proxeus/sports/2.png)‍  

## Implementation


An athlete receives a link, fills out the forms and finalizes the process before sharing the finalized document with a doctor (via the address book) and with the medical board of ITA in charge of issuing the final confirmation.  

Finalized TUE documents can be easily verified on a drag and drop interface, which triggers the creation of a new hash. The system checks on the blockchain whether the exact same hash has been registered, allowing to prove whether the document is original or has been tampered with.  

The process was presented on CNN Money Switzerland where Benjamin Cohen, director of ITA, explained the genesis of the organization, its challenges and how blockchain could play a role in the future.  

![](_media/old_proxeus/sports/3.png)

## Insights


The project was implemented as a prototype through an active exchange with ITA and gave the team an occasion to reflect on the security aspects of the application, as well as to add, test and improve several features, such as the signature requests, the rendering of tables in documents and the possibility to upload custom images as part of a workflow, which can be displayed as attachments in the final document.

## Limitations


The lack of easy access to blockchain applications (e.g. integration of blockchain identity solutions in smartphones) remains an issue for products to be implemented at wide scale. In the present case for example, athletes cannot be expected to perform blockchain transactions independently.

## Try it out yourself


You enjoyed reading the documentation and would like to build a similar workflow? Follow the steps below. All you need is an instance of Proxeus and basic users skills. Check out the handbook and tutorial videos on our website for help.  
  
Here is how you can create a workflow - using only Proxeus and no programming at all.  

*   Understand the requirements. What workflow output do you expect? What documents should be registered? What role should signatures play?
*   Set up Proxeus. You can run your own instance of Proxeus on a server or locally on your computer - or you could use someone else’s instance. The complete guide to setting up your own instance is available [here](http://doc.proxeus.com/#/README). It is recommended that you deploy your own smart contract following our instructions in the guide and using the [template in our GitHub](https://github.com/ProxeusApp/proxeus-contract).  
    
*   Prepare the document template(s). What should the document design look like? What should placeholders be used for? How should the information be formatted (e.g. what sections should be shown vs. hidden in the final output)?  
    
*   Prepare the user forms. Adapt the requirements to the audience (e.g. athletes cannot be expected to be very tech-savvy) - entering the data has to be intuitive and instill confidence. All Proxeus form elements support help texts - they even accept HTML and links to further information - use them to clarify every step.
*   Configure the workflow by connecting the smart template(s) and the forms in a workflow. Set a XES price to the workflow (if applicable) and share the form with the platform users.
*   Set signatures as required - once given by users they are publicly visible on the blockchain and verifiable by anyone.
*   Test and improve your workflow. Nothing is perfect from the start. Even when your workflow is already in production and in use, you can simply clone it and release an improved version to your users. If the changes are compatible, you can also just upgrade the existing workflow and all running instances will use the new definition automatically.
*   Create an instruction page to onboard users. Our handbook and tutorial videos are also here to help.
*   Integrate a verification tool into your organization’s website to simplify the verification process. The verification tool of your Proxeus platform can be found directly via the /validation path; for example [https://morrison.proxeus.org/](https://morrison.proxeus.org/). It can also be embedded using a standard iFrame HTML code.

  
Your iFrame HTML code to embed it into your website would then look like this:

<iframe width="100%" height="650" src="https://morrison.proxeus.org/" frameborder="0" marginwidth="0" marginheight="0" scrolling="yes"></iframe>

Tokenization
============

**Tokenization of Luxury Cars**

![](_media/old_proxeus/tokenization/1.svg)

## Background


Tokenization of assets is a hot topic and a collaboration was set up in 2018 with Mercuria Helvetica, a startup playing a pioneering role looking at ways to tokenize investment-grade cars, to explore the possible role of the workflow and document automation capabilities of Proxeus in that context.

![](_media/old_proxeus/tokenization/2.png)


## Use case exploration


The tokenization of an asset means its representation in the digital world through a digital artifact (e.g. using a hash or a token on a blockchain). The creation of multiple tokens each representing a piece of an asset can make this one tradable if the tokens are sold and put onto a marketplace.

While there are several ways and methods to reach this result, we defined and followed the following steps:  

1.  Documentation - define the asset you want to tokenize. Build the documentation you need to describe fully all the aspects of your asset that make it unique, preparing it for registration.
2.  Registration - immutably register your asset. Once your asset is fully documented, register it on the blockchain in order to create a tamper-proof record.
3.  Certification - bring in third-party experts to certify the accuracy and completeness of your documentation, possibly after inspecting the real-world asset. Let them add their cryptographic signatures to your document as evidence for their approval.  
    ‍

If fractional ownership is a topic, the issuer can decide to generate tokens by deploying a smart contract for the asset described in the documentation as an additional step.

## Implementation


The first step was to identify everything that needs to be documented about the luxury cars to be tokenized. Besides the technical attributes of the car (engine type, horse power etc.), we learned that investors will request to know price-defining data such as the production year, rarity (i.e. number of cars produced for a certain model), the history of repairs and maintenance performed, the car’s current condition and so on. The technical description is typically accompanied by a photographic documentation.  

Based on these requirements, the team designed user-friendly data entry forms. The forms were configured to validate the input immediately (e.g. valid date, value within the expected range, all mandatory fields filled etc.)

  

![](_media/old_proxeus/tokenization/3.png)

_An excerpt from a user form draft in the form editor_  

‍

![](_media/old_proxeus/tokenization/4.png)

_The configuration for an input, validated to be a number and bigger than 0_  

  

As the second step the team designed a smart template that presented the car’s information as well as the uploaded photos in an appealing manner. The template would show or hide sections depending on the information entered in the workflow.

![](_media/old_proxeus/tokenization/5.png)

_An excerpt from the document template (draft version)_  

Thirdly, the forms and template were built into a workflow designed to be used by the Mercuria Helvetica’s staff. It allowed for the efficient entry of all relevant data, the rendering of the documentation into a PDF file as well as an effortless registration of the produced documents in a Proxeus smart contract on the Ethereum blockchain.

![](_media/old_proxeus/tokenization/6.png)

_A workflow draft that connects four user forms and results in the production of a blockchain-secured documentation of a luxury car_

## Result


The team deployed a new instance of the Proxeus platform and designed the complete workflow to create the documentation dossier. The workflow included the blockchain registration. With this application, Mercuria Helvetica received a standardized and streamlined process to create the documentation for each car to be tokenized. It also provided a tool that enabled anyone to quickly verify the authenticity of all circulating documents.  

Our partner in this project later went on to create Curio Invest, a platform that enables the issuance, distribution and management of fine collectable assets registered as security tokens. CurioInvest makes it possible for collectible assets such as fine automobiles to be registered with the regulators in a compliant way and to be represented as a token that users can acquire.  

## Feedback


In an [interview](https://medium.com/proxeus/blockchain-enables-investing-in-classic-cars-for-everybody-not-just-millionaires-5e9dec740aee) from August 2019, CurioInvest’s CEO Fernando Verboonen made the following statements:  

_“Last year we did a proof of concept using blockchain to register certificates of authenticity and conditions of investment-grade classic and exotic autos. By scanning those certificates and placing the document hash on the blockchain, we were making them verifiable, guaranteeing their originality and authenticity._  

_By having every document involved in the acquisition of a classic car immutably on the blockchain, you are able to scale the process, making it safe, decreasing diligence costs and saving a lot of time: we’re reducing the process from several weeks to seconds!_  

_Our current platform is focusing on the tokenization of the assets and is still quite minimalistic, but in the future we see the opportunity of using Proxeus to create certified documents that can be referenced in the relevant tokenization smart contract. This step will further enhance CurioInvest’s scope and usability by adding features like integrating all the relevant documents of the classic cars.”_

## Insights


The project allowed the product team to understand the possible role of a document-centered workflow engine like Proxeus in the tokenization space, differentiating between the processes that should take place within the application core (document creation, registration, certification) and the ones that are custom to every project and should be set up separately and under the responsibility of the respective projects (smart contract deployment, custody solutions).  

## Limitations


The discussions around the requirements for certifying assets represented critical inputs for designing the Signature Requests feature, which was not available at the time, but now enables certification with Proxeus out of the box.

## Try it out yourself


If you’ve enjoyed reading the documentation of this project and would like to try building something similar, we suggest following the steps below. All you need is an instance of Proxeus and decent users skills. If you haven’t already, you should check out our extensive [handbook](https://docs.google.com/document/d/1C3B1oNY6lOv8Q_AvbKhwlySrS6qTiRl3raPLV6OXr7w/preview) and our [step-by-step guide](https://docs.google.com/document/d/1Gl6R1t0LYRK6kARScx5-vqpiKtd6xqcH2yPPC3wWMr4/preview) for creating a workflow.

‍

1.  Understand the requirements. What workflow output do you expect? What documents should be registered? What role should signatures play?

  

2.  Set up Proxeus. You can run your own instance of Proxeus on a server or locally on your computer - or you could use someone else’s instance. The complete guide to setting up your own instance is available [here](http://doc.proxeus.com/#/README). It is recommended that you deploy your own smart contract following our instructions in the guide and using the [template in our GitHub](https://github.com/ProxeusApp/proxeus-contract).  
      
    
3.  Create the template(s) for the document(s) you want to produce, register and get certified. What are the properties of the asset that need to be described in its documentation? What do the future buyers of the tokens need to know? What are the design requirements?  
      
    
4.  Create the user forms. Now that you’ve designed the desired result, you know what information needs to be collected to fill the templates. Use the data validation capabilities of Proxeus to help the user avoid errors.

  

5.  Configure the workflow by connecting the smart template(s) and the forms in a workflow. Set a XES price to the workflow (if applicable) and share the form with the platform users.

  

6.  Test and improve your workflow. Nothing is perfect from the start. Even when your workflow is already in production and in use, you can simply clone it and release an improved version to your users. If the changes are compatible, you can also just upgrade the existing workflow and all running instances will use the new definition automatically.  
      
    
7.  Create an instruction page to onboard users. Our handbook and tutorial videos are also here to help.

  

8.  Set signatures as required - once given by users they are publicly visible on the blockchain and verifiable by anyone.  
    ‍

‍  

Integrate a verification tool into your organization’s website to simplify the verification process. The verification tool of your Proxeus platform can be found directly via the /validation path; for example [https://morrison.proxeus.org/](https://morrison.proxeus.org/). It can also be embedded using a standard iFrame HTML code. 
  
Your iFrame HTML code to embed it into your website would then look like this:

<iframe width="100%" height="650" src="https://morrison.proxeus.org/" frameborder="0" marginwidth="0" marginheight="0" scrolling="yes"></iframe>


## Take it one step further


Deploy a smart contract!  

After picking an asset, compiling extensive documentation, anchoring it on the blockchain and getting it certified by experts, you may want to deploy a token smart contract for this asset. There is no universal smart contract solution and you should take time to consider different aspects to find the one that fits your needs: how many tokens should be issued? Should they be tradeable? Who will issue them? Should they be burnable? By whom? Should the supply be fixed or should further minting be possible? Is the smart contract controlled by a single private key or do changes require the signature from multiple keys?  

In most jurisdictions, tokens may be considered securities. The consultation of legal experts is necessary if you plan to distribute the tokens created in any way. On the technical side, the deployment of smart contracts entails critical security aspects and this task should be given to specialists, in particular if your tokens are planned to have monetary value (for example by representing a physical asset): smart contract security breaches with partial or total loss of funds have happened frequently in the past and in the famous “DAO hack”, 50 million dollars worth of ETH were lost due to a loophole.  

An example contract template for the registration and certification of documents (“ProxeusFS”) is available in our GitHub repository. For token contract resources, we recommend to look into [OpenZeppelin](https://github.com/OpenZeppelin/openzeppelin-contracts/tree/master/contracts/token/ERC20) open source templates. If you are looking for support, tokenization services are offered by many companies, including [Microsoft Azure](https://azure.microsoft.com/en-us/services/blockchain-tokens/), [Token Factory](https://tokenfactory.global/), [Tokeny](https://tokeny.com/) and [AlphaPoint](https://alphapoint.com/) to name just a few.  

There are lots of free and helpful resources out there. And feel free to reach out to the Proxeus community for help with your endeavors!

Logistics
=========

**Track and Trace Prototype**

![](_media/old_proxeus/logistics/1.svg)

## Abstract

On this page we’re presenting one of our exploratory projects that helped shape the development of Proxeus. It demonstrates the use of the platform in a logistics use case. We’re tracking shipped goods from their origin all the way to their destination using RFID tags, the [IOTA](https://www.iota.org/get-started/what-is-iota) ledger and Proxeus. Together with a description of the scenario and our solution you’ll also find a guide for the reproduction of our project.

## Background

IBM has estimated that the implementation of digital technologies such as blockchain could save the logistics industry as much as $38 billion per year. According to Boston Consulting Group, blockchain technologies alone could save several millions of dollars yearly for large manufacturers.

Supply chain processes traditionally rely on paper, for example a full binder of documents is needed to send a shipment from Asia to Europe. A key saving, in financial terms, would come from a reduction of the manual paperwork associated with the movement of goods as they transit in between terminals.

This led us to define “Logistics & Trade Finance” as one of the focus verticals that we set out to explore as part of the Proxeus project development. Indeed by collaborating with external partners on concrete use cases, we were able to gain valuable insights that contributed to shape the Proxeus framework.

Today, the Proxeus document workflow automation, blockchain capabilities and integration interfaces provide a solid foundation for projects in logistics and trade finance.

## Use Case Exploration

By combining new technologies such as the Internet of Things (IoT) together with a decentralised ledger, we looked into possibilities to track shipments through the use of physical devices while registering the metadata associated with a particular shipment in a transaction on a ledger, making possible to instantaneously generate digital documents providing reliable and immutable tracking information for the transit of goods from A to B.

## What we've built

The prototype supported a full showcase process using a real logistics industry use case: the tracking of goods from the seller to the customer, going through one or several carriers on the way, coordinated by a freight forwarding company.

<figure>
    <iframe allowfullscreen="true" frameborder="0" scrolling="no"
            src="https://www.youtube.com/embed/hburhJQy3m8">
    </iframe>
</figure>

Watch our video explaining the track & trace process

## Actors

These are the key participants in the freight forwarding process.

Shipper: Seller of the goods, engaging a Freight Forwarder to organize transport of the goods from point of origin to the agreed destination with the customer.

Freight Forwarder: Service Provider to the Shipper organizing transport from point of origin to destination incl. hand-over to the customer.

Carrier: Vessel operator contracted and coordinated by the Freight Forwarder to transport the goods for the full or part of the distance between point of origin and destination.

Customer: Buyer of the goods expecting delivery at destination.

![](_media/old_proxeus/logistics/2.png)

An example visualization of the goods and information flows in logistics.

‍

## The process

We’ve selected a (slightly simplified) freight forwarding process for our project. These are the key steps:

- The Seller sells goods to the Customer and engages a Freight Forwarder to organize transport from point of origin to destination.
- The Freight Forwarder responds to the Shipper’s request by issuing a shipping quote which is valid for a certain timeframe. Upon quote acceptance by the Shipper, the booking confirmation is added to the shipment collection which has a unique identifier, i.e., via the Global Identification Number for Consignment (GINC).
- The Shipper packages the goods on to the pallets, and assigns the pallets to the GINC. The chips are activated for that shipment and set into status “Ready for pick-up”, which also marks the first step of the tracking flow on the blockchain.
- Out of his system data, the Shipper produces the following items that are added to the shipment collection:

  a. Commercial Invoice for shipment  
  b. Packing List (overall and by pallet)  
  c. Load Details for Bill of Lading (shared with Freight Forwarder)

- The truck carrier picks up the shipment paletes at the point of origin and checks them in on his vessel (HT), which issues the official Bill of Lading in the shipment collection and makes it available to the Shipper and puts the pallet as in transit leg 1 of 2, the first update of the blockchain tracking flow.
- At the reloading terminal, the HT in the forklift checks in the receipt of every pallet and the status of the shipment is changed into reloading leg 1 to leg 2, the second update of the tracking flow.
- Upon loading by the train carrier, the pallet is checked in to leg 2 of 2, the third update of the tracking flow.
- At arrival at destination, C checks in the pallet (fourth and last update of the tracking flow) and accesses his relevant documents to confirm receipt (finalization of tracking flow). Pallets (and therefore their chips) are deactivated for shipment.

## Insights

We’ve learned quite a few interesting things during the course of this project. For example, who would have thought that it may become a project requirement to run Proxeus not on today’s powerful cloud infrastructures with endless performance scalability, but on a tiny computer like the Raspberry Pi? This move reinforced the team in architectural decisions such as splitting the document service from the Proxeus core and enabling its use as a standalone service.

The project provided a lot of insights for the product management and development team and helped improve the API capabilities of Proxeus - through the Proxeus API, even microcomputers like Raspberry Pis are able to trigger the production of a whole series of different documents - from anywhere in the world.

Reading data from an [IOTA tangle](https://www.iota.org/) - a technology very different from Proxeus’ native blockchain Ethereum - and feeding this data into Proxeus workflows for the production of documents worked really well. For the team it was important to see that it works not only in theory - you can actually connect Proxeus to any blockchain thanks to its extensibility.

## Limitations

The prototype was intentionally created for demonstration and learning purposes only and therefore lacks many elements that would be required of a full-fledged logistics solution. There were no integrations with the existing systems of any stakeholder. Only prototyping hardware was used instead of devices actually hardened for the use in logistics. Everything was simulated in one room instead of taking place all over the world with all the implications of different time zones, climates, procedures and so on.

The actual implementation of a blockchain solution into a vast existing ecosystem such as global logistics will require to involve and coordinate the needs of a great many different stakeholders. Even if a lean scope is chosen, such a project is a great endeavour and will require substantial initial investments from all involved parties (hardware, integration projects, employee training etc.).

## How we’ve built it

Let’s dive into the technical details! Make sure you consult the glossary at the end of the page if you encounter any unfamiliar terms.

As part of a project executed in cooperation with a German logistics startup, RFID chips were used to track palettes across connected terminals. The tracking data (such as the consignment ID, the terminal name or even the signatory name) was directly sent to the IOTA ledger by the terminals without any human intervention whatsoever.

On a separate server, Proxeus was then leveraged to automatically and instantaneously generate documents based on the data retrieved from the IOTA ledger and pre-configured templates. Several key documents (such as the [Bill of Lading](https://en.wikipedia.org/wiki/Bill_of_lading)) can be tokenized to enable a frictionless transfer of ownership of the goods shipped.

![](_media/old_proxeus/logistics/3.png)

Image: palette with goods and an hidden RFID chip being “transported” to a terminal

![](_media/old_proxeus/logistics/4.png)

Image: screen with custom-made demo “terminal” software running on a Raspberry Pi

‍

In this prototype, we introduced quite a few elements in order to showcase a potential cross-platform, blockchain-agnostic tracking system that could be built using Proxeus and IoT components.

## IoT node on a Raspberry Pi**

**‍**In the world of the “internet of things (IoT)”, the participating devices are often called “nodes”. For the nodes in our simulated logistics ecosystem, we equipped a prototype board (Raspberry Pi) with an RC522 RFID reader. The microcomputer connects to WLAN automatically and the operating system is stored on an SD card which can easily be cloned after setting it once.

We created a small custom application for the Raspberry Pi. Running this software, the nodes were able to run independently in different modes:

1\. Track consignment

2\. Sign receipt for a consignment

3\. Open up a new tracking

4\. Set terminal in monitoring mode

5\. Display tracking details

For the technically curious among our readers, the code for the tracking nodes can be found on our GitHub: [https://github.com/ProxeusApp/usecase-shipment-tracking](https://github.com/ProxeusApp/usecase-shipment-tracking)

Keep in mind that it was a rapidly built prototype and may need some improvement if you wish to use it for your own projects.

One part is built in Python and offers functionality to read from and write to RFID tags using the [SimpleMFRC522](https://github.com/pimylifeup/MFRC522-python) library specifically made for budget MFRC522 module for the Raspberry. It also creates and tracks transactions on the IOTA ledger using the respective [IOTA library for Python](https://github.com/iotaledger/iota.py). It has a console-based user interface with a menu to choose between the modes listed above. The second part is a visual UI implemented in Golang that can be used on top.

After the initial setup, the application was configured by connecting to it and copying the specific configuration files to it. In order to provide maximal mobility, the microcomputer was connected to a power bank. This setup is barely larger than a cigarette box (depending on the powerbank chosen) and is connected to the internet over WLAN. The platform also incorporates an IOTA node which syncs the data logged by the application with the livenet of IOTA.

The IoT nodes in this project came in two variants:

1.  Priming nodes are writing information to the RFID chip, linking it to a unique id like a dossier number. The second information is the UID of the tag itself, which cannot be changed and identifies the pallet.
2.  Tracking nodes are always on “listen mode” after bootup and execute an IOTA transaction containing the details (tagID, document collection, data & time etc) of the tracking.

![](_media/old_proxeus/logistics/5.png)

Image: a Raspberry Pi with an RFID module (blue), powered by a power bank

As soon as an RFID chip is scanned, the terminals execute a transaction. The software also returns a table with all the checkpoints encountered. The system has been configured to be able to run concurrent trackings and have them added to the respective IOTA tangle.

## Documents as a Service‍

In this use case, Proxeus was used as a “documents as a service” provider. Proxeus can be remotely controlled to produce documents by using the API. Any scripting or programming language that is capable of communicating with a REST API will do. This could also be simulated by using a tool like [SoapUI](https://www.soapui.org/).

1.  To make the call, you must be [authenticated](http://doc.proxeus.com/#/api_auth). You can create API keys in your user profile on the Proxeus platform.
2.  Write down your workflow’s ID or retrieve it via the API’s [workflow list request](http://doc.proxeus.com/#/api_list_all_workflows).  
    GET /api/document/list  
    For the example, let’s assume it’s 3e6ece3d-6b5d-4e79-aea0-0c06e14935cb.
3.  Retrieve the workflow’s [data schema](http://doc.proxeus.com/#/api_get_workflow_schema) to identify the required inputs.  
    GET /api/document/3e6ece3d-6b5d-4e79-aea0-0c06e14935cb/allAtOnce/schema
4.  Call the API method “allatOnce” as documented [here](http://doc.proxeus.com/#/api_execute_workflow).  
    POST /api/document/3e6ece3d-6b5d-4e79-aea0-0c06e14935cb/allAtOnce  
    {“field”:”value”, “field2”:123, “field3”:true}  
    The success response is status 200. The response body will be the produced PDF.

The advantage of a REST API is that you’re free to use almost any technology. You could write a little script in Python or use programming languages like Go or Java.

A Proxeus platform we’ve set up specifically for this project was equipped with a number of workflows with templates for the relevant documents (shipping manifest, bill of lading, container load plan, confirmation letter, tax document, import and customs documents). For this we researched what information these documents typically contain and designed the respective templates in Microsoft Word (OpenOffice works as well). In the Proxeus admin panel we created one workflow for every template. Then we configured simple forms (mostly text/number/date inputs) for the required inputs and added them to the workflows. All this was done using standard Proxeus functionality - if you’re not familiar with it, we highly recommend watching our \[tutorials\] or browsing the \[handbook\]!

![](_media/old_proxeus/logistics/6.png)

Screenshot: an excerpt from the bill of lading template we created

‍

## Monitoring service

A key element of the showcase was to demonstrate that the transactions have indeed been logged on the blockchain and are made available to the user. This has been achieved by using an IOTA tangle explorer. The monitoring service can be used to simply display updates on a screen or to trigger document production via the API of a Proxeus instance. We’ve described in the chapter above how to use Proxeus in a “documents as a service” setup.

A transaction on the IOTA tangle as viewed through the popular tangle explorer [TheTangle.org](https://thetangle.org/) looks like this:

![](_media/old_proxeus/logistics/7.png)

## How to build it yourself

Are you interested in building a similar project? The following steps can guide you through the process. Please note that this is one of the most challenging Proxeus projects to rebuild. It requires prototyping hardware and you will need **software development skills** as well as a **good level of experience in using Proxeus**. Read our solution description above carefully before following this guide.

1.  If you haven’t done this already, please familiarize yourself with Proxeus workflows.The documentation on the website and in particular the following materials should help:  
    a) \[link to handbook\]  
    b) [Tutorial video - how to build your first workflow](https://youtu.be/8ndmlTPYwMc)  
    c) [Tutorial video - Create and manage documents](https://youtu.be/WwyOiluZw8c)  
    d) [Step by Step Guide - Education Diploma](https://docs.google.com/document/d/1Gl6R1t0LYRK6kARScx5-vqpiKtd6xqcH2yPPC3wWMr4/preview)
2.  Requirements analysis and specification:  
    a) Which part of which business process would you like to digitize with your Proxeus prototype?  
    b) Who are the stakeholders? What are the different stations of the transiting goods? What are their roles? Do they simply report that the goods have passed through or do they need to produce a receipt for someone?  
    c) Which documents shall be produced and what information is needed to create them?
3.  Preparation of materials  
    a) [Raspberry Pis](https://www.raspberrypi.org/) (or [Arduinos](https://www.arduino.cc/))  
    b) SD memory cards for the Pis  
    c) Battery packs for the Pis  
    d) Screens for the Pis  
    e) RFID readers/writer modules for the Pis, RC522 standard  
    f) RFID chips or cards, RC522 standard  
    g) Optional: to visualize the transportation of goods, we’ve used props like miniature pallets, self-built miniature trucks and boats. Add an RFID chip to each pallet.
4.  Prepare an RFID priming node: one of the Raspberry Pis with an RFID writer module can do this job. It should create a unique ID for each new tag that it finds in its proximity. It can be implemented in the scripting/programming language of your choice, as long as it is supported by one of the available operating systems of the device. [This tutorial](https://lastminuteengineers.com/how-rfid-works-rc522-arduino-tutorial/) is a great starting introduction if you haven’t worked with RFID before and this [detailed guide](https://pimylifeup.com/raspberry-pi-rfid-rc522/) walks you through all the Raspberry specific steps.
    Our implementation for Raspberry Pi can be found on our GitHub:  
    [https://github.com/ProxeusApp/usecase-shipment-tracking  
    ](https://github.com/ProxeusApp/usecase-shipment-tracking)It triggers a transaction on the IOTA tangle when an RFID chip is detected by the Raspberry’s attached RFID module.
5.  Set up the listening nodes- the Raspberry Pis with RFID reader module - with their specific tasks depending on the specifics of your desired setup. Just like with point 4 you are free to choose the technology for this script or application. This could be simply reporting that the goods have passed this station or they could display information. They could also trigger the production of a receipt or of a payment. When an RFID chip is scanned, the node should automatically create a transaction on the IOTA tangle with the ID of the tag and the node’s identification. For this, the Raspberry Pi needs to be connected to the internet. The node’s software must be capable of creating IOTA transactions (e.g. through the API of an IOTA service provider). A simple script (e.g. Python) should suffice. Here or  
    \- [Python quickstart guide to IOTA  
    ](https://docs.iota.org/docs/client-libraries/0.1/getting-started/python-quickstart)\- [PyOTA Getting Started  
    ](https://pyota.readthedocs.io/en/latest/getting_started.html#getting-started)\- [How to create a IOTA transaction in a single API call (Python)  
    ](https://gist.github.com/Hribek25/abf3e51864c32d2e0df5e20785d0cb36)\- [A popular Tangle explorer for IOTA](https://thetangle.org/)  
    \- You can use the following public nodes service: [https://nodes.thetangle.org:443

    ](https://nodes.thetangle.org:443/)

6.  Set up Proxeus. You can run your own instance of Proxeus on a server or locally on your computer - or you could use someone else’s instance. The complete guide to setting up your own instance is available [here](http://doc.proxeus.com/#/README) and in our GitHub repository.
7.  Create Proxeus workflows. When you scoped your project in step 1, you figured out which documents need to be produced. After you’ve set up your own instance of Proxeus, you can now configure a workflow for each document. Workflows comprise of data entry forms and document templates. As you’re not going to type in the data manually, you don’t need to spend much time on designing beautiful forms. Adding all necessary fields with the right data type should suffice. Design and upload templates for all document types and add them to your workflows.  
    \- \[Download link for template - Bill of Lading\]  
    \- \[Download link for template - Consignment Tracking History\]

8.  Set up the event listener. On your computer, create a software (e.g. a Python script) that listens to the events (transactions) that your Raspberry Pis (the listening nodes) create on the IOTA tangle and then triggers a Proxeus workflow through its API. The API allows you to choose the workflow for the type of document you want to produce and to provide all necessary data to immediately complete it and generate the document. For example: if the final node in your logistics chain scans the pallet, it creates a transaction that is then detected by your listener. It then triggers Proxeus to produce a confirmation for receiving the goods. Examples for using the Proxeus API can be found above under “Documents as a Service” and the documentation is at [doc.proxeus.com](http://doc.proxeus.com/#/api_auth).

‍

## Glossary

We’re using a lot of terms that you may not yet be familiar with. Here is some help.

IoT

The Internet of Things is a network of devices that are interacting with each other without human interaction. [Read more on Wikipedia](https://en.wikipedia.org/wiki/Internet_of_things)

IOTA

IOTA is a distributed ledger technology (DLT) that allows devices in an IOTA network to transact in micropayments and to send each other immutable data. [Visit the website](https://docs.iota.org/)

Tangle

IOTA uses the term “tangle” for what is called the blockchain in other distributed ledger technologies: an immutable data structure that contains all transactions.

TheTangle.org

A service provider offering access to the IOTA tangle through a browser on their website and an API. It is much like what Etherscan is for Ethereum. [Visit the website](https://thetangle.org/)

RFID

A near-field data transmission technology. [Read more on Wikipedia](https://en.wikipedia.org/wiki/Radio-frequency_identification)

Raspberry Pi

A tiny, simple computer developed for teaching the basics of computer science. Popular for technical experiments and prototyping. [Visit the website](https://www.raspberrypi.org/)

API

A programming interface to facilitate the communication between different computer programs.

Python

A powerful scripting language. [Visit the website](https://www.python.org/)

SoapUI

An advanced API testing tool that enables you to try out any API without writing a single line of software. [Visit the website](https://www.soapui.org/)

## Take it one step further

As we have mentioned under “Limitations”, this project was limited to a tightly defined scope. There are many more directions to explore.

Decentralized IoT networks may transform supply chain management in the future. Most of today’s processes involve the management of physical, hand written documents. All those formal paper-based documents and processes need first to be digitized and later integrated with other systems.

Using Bluetooth sensors combined with other IoT devices would enable users to retrieve and add additional data such as temperature, pressure or any other quantifiable variables to the tracking journey, opening a new dimension and new possibilities for implementing documents (e.g. contract execution if a certain condition is met).

Projects like [Ambrosus](https://ambrosus.com/), [Modum](https://modum.io/) and Zimt have specialized in this. Read our interview with Vlad Trifa, a thought leader in the fields of IoT and blockchain [here](https://medium.com/proxeus/our-mission-is-to-provide-trusted-real-time-data-as-a-commodity-aeb5a1d5118d).

Data Storage
============

**Proxeus Releases Decentralized Data Storage DApp Tango Reloaded**


April 11, 2019

If you have followed us for a while, you will already be familiar with the DApp Tango Beta. Tango is a user-friendly, ground-breaking secure data storage product which empowers anyone to register, share, sign and save documents through a decentralized blockchain application.

![](_media/old_proxeus/data_storage/1.png)

We are now excited to announce the release of Tango Reloaded, a major update incorporating new features and improvements following extensive user testing over the last months. The difference between Tango Reloaded and traditional storage solutions, such as Google Drive or Dropbox, is that Proxeus puts the power entirely in the hands of the user.

How? The ownership and sharing rights are held by the user’s crypto-identity, the files are fully encrypted and the storage provider can be selected individually for each file based on the physical location or security grade required by data protection requirements.

## Your data as secure as Fort Knox


Proxeus stores the encrypted data in a secure off-chain facility, while putting advanced access controls on-chain using smart contracts.

ROCKCHAIN, the service provided by [Mount 10](https://www.mount10.ch/), is currently the exclusive storage provider, storing the data in what can best be described as a military bunker, [Fort Knox style](https://www.mount10.ch/de/mount10/swiss-fort-knox/index.html). It is Europe’s most secure data centre, with multi-level protection and servers stored deep inside the Swiss mountains.

Co-founder Antoine Verdon says: “I am very excited about Tango as its potential goes way beyond the decentralized storage of information. We allow users to truly own their data, to associate it with their crypto-identity and to control the sharing of it over a blockchain interface. It is a key piece of technology that constitutes a stepping stone for anyone interested to build blockchain products around the concept of self-sovereign identity.”

## Powerful, useful features


The new features offer enhanced functionalities and have been defined in collaboration with several community members who gave feedback in the months following the first Tango release.

### Key features include:


*   Zero Knowledge Storage (achieved through client-side PGP encryption), meaning no one (Including the storage provider) has access to the file’s content, unless the user decides to share the file
*   Blockchain secured access management (including sharing).
*   Convenient document signing with your private key.
*   DApp Tango operates as a “pay-as-you-go” solution, meaning there is no need for a traditional subscription. The payment for individual file storage is made with Proxeus’ own cryptocurrency, XES, and is now based on the file size and the duration of storage.
*   Users can send Ether and XES straight out of their DApp wallet.
*   A notification log, so users can track activity in one place.
*   Minor bug fixes and general improvements.
*   Local Encryption Concept which improves file handling on the user’s local machine.

## Extremely user-friendly


Thanks to new ground-breaking advancements, now literally anyone can store files and manage permissions decentrally. “Tango Reloaded is incredibly simple to use. The interface is extremely user-friendly, so anyone can use it to upload and share, sign and register documents, says Stephan Balzer, Product Manager at [BlockFactory AG](https://www.blockfactory.com/), the company that has developed the DApp for Proxeus.

## Why the name Tango?


The Proxeus team seems to have an affection for passionate dances. The name Tango was chosen as the dance represents strength and firmness, but simultaneously, a vibrant and playful energy. And you need both elements to make the technique work.

## Test Tango Reloaded


Curious? Test Tango Reloaded [here](http://beta.proxeus.com./), and check out [this video](https://www.youtube.com/watch?v=B7bgXY2yGGw&feature=youtu.be) that walks you through the simple steps to use it. There’s also a user-guide available just in case you get stuck.

When you’ve tested it, we want to hear from you! What are your favourite features? Would you like different functionalities and if so, what? [Let us know.](http://proxeus.com/)

‍