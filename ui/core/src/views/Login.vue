<template>
<div>
  <!--<sign-message id="sign-message"></sign-message>-->
  <vue-headful :title="$t('Proxeus - Sign in','Proxeus - Log in')"/>
  <h1 class="text-center">{{$t('Sign in title','Log in')}}</h1>
  <div class="login-form container-fluid px-4 pt-2 mt-3 bg-light" :class="{'login-form-sm': !metamaskLoginAvailable}">
    <div class="row">
      <div class="col-12 mt-3 text-center d-flex align-items-center" v-if="app.wallet && metamaskLoginAvailable">
        <div class="mid align-self-center w-100" v-show="signing === false">
          <h2 class="mb-3 font-weight-bold">{{$t('Account signature')}}</h2>
          <p class="light-text">{{$t('Account to sign in explanation','Use your MetaMask account to log in.')}}</p>
          <p class="text-danger" v-if="walletErrorMessage">{{ walletErrorMessage }}</p>
          <button class="btn btn-primary px-3" @click="metamaskLogin"
                  v-if="metamaskLoginAvailable">{{$t('Sign in with signature','Log in with signature')}}
          </button>
        </div>
        <div class="mid align-self-center w-100" v-show="signing">
          <spinner background="transparent" style="position:relative;" :margin="60"></spinner>
          <p
            class="text-primary">{{$t('Account signature explanation', 'Please sign with MetaMask. If MetaMask doesn\'t pop up automatically, click on the MetaMask icon in the top right corner of your browser.')}}</p>
        </div>
      </div>
    </div>
  </div>
  <div class="modal fade" ref="tcModal" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-lg" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Proxeus Terms and Conditions of Service</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body px-3">
          <h2>Proxeus Terms and Conditions of Service</h2>
          <div><p class="c15"><span class="c0"></span></p></div>
          <ol class="c8 lst-kix_list_51-0 start" start="1">
            <li class="c7"><h1 style="display:inline"><span class="c6 c12">General</span></h1></li>
          </ol>
          <p class="c1"><span class="c4">These general terms and conditions (hereinafter the &ldquo;</span><span
            class="c14">Terms</span><span
            class="c4">&rdquo;) set out the rights and obligations in connection with the use of </span><span
            class="c4">the services </span><span class="c4">offered by Proxeus Anstalt, in Vaduz, </span><span
            class="c4">Liechtenstein</span><span class="c4">&nbsp;(hereinafter &ldquo;</span><span
            class="c14">Proxeus</span><span
            class="c4">&rdquo;, &ldquo;</span><span class="c14">we</span><span
            class="c4">&rdquo;) and the user (hereinafter &ldquo;</span><span
            class="c14">User</span><span class="c4">&rdquo;, &ldquo;</span><span class="c14">you</span><span
            class="c4">&rdquo;</span><span
            class="c4">) in connection with the utilization of the Services, as well as further services, applications and functions that are offered by Proxeus through its platform, unless explicitly stated otherwise.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1 start" start="1">
            <li class="c11"><h2 style="display:inline"><span class="c2">Usage</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">These Terms shall become legally binding and enforceable upon your first using the Services as defined below in section 3. </span>
          </p>
          <p class="c1"><span
            class="c4">The User has the opportunity to print these Terms (representing the text of the contract created) at any given time. For this purpose the print function of the respective browser can be used.</span>
          </p>
          <p class="c1"><span
            class="c4">The User has no claim to the conclusion of a contract. Proxeus reserves the right to refuse the offer of a User to enter into a contract with them, without giving reasons.</span>
          </p>
          <p class="c1"><span
            class="c4">Users of the services offered by Proxeus must be at least 16 years old. Proxeus preserves the right to review compliance with this requirement any time after registration. If we become aware that we have collected personal data of persons aged less than 16, we shall delete this data immediately; unless we are legally obligated to retain the data.</span>
          </p>
          <ol class="c8 lst-kix_list_51-0" start="2">
            <li class="c10"><h1 style="display:inline"><span class="c6 c12">Account</span></h1></li>
          </ol>
          <ol class="c8 lst-kix_list_51-1 start" start="1">
            <li class="c11"><h2 style="display:inline"><span class="c2">Account password and security</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">When setting up an account with Proxeus, you will be responsible for keeping your own account information and keeping it secret. You agree to (a) never use the same password for the Services as defined below that you have used outside of the Services; (b) keep your account information and password confidential and do not share them with anyone else; (c) immediately notify Proxeus of any unauthorized use of your account or breach of security.</span>
          </p>
          <h2 class="c1" id="h.gjdgxs">2.2 Termination and suspension of your account
          </h2>
          <p class="c1"><span
            class="c4">We reserve the right to suspend the User&rsquo;s access to the Services as defined below either temporarily or permanently if there are concrete indications that the User has violated or shall violate these Terms and/or the law, or if Proxeus has a legitimate interest in suspending the User&rsquo;s access.</span>
          </p>
          <p class="c1"><span
            class="c4">In deciding as to whether access to a User shall be suspended or terminated, the legitimate interests of all parties will be considered as appropriate.</span>
          </p>
          <ol class="c8 lst-kix_list_51-0" start="3">
            <li class="c7"><h1 style="display:inline"><span class="c6 c12">The Services</span></h1></li>
          </ol>
          <p class="c1"><span
            class="c4">Proxeus provides its Users with online services through its platform to build blockchain applications amongst other things (the &ldquo;</span><span
            class="c14">Services</span><span
            class="c4">&rdquo;). The Services are provided in a hosted environment and interact with a distributed ledger technology (hereinafter &ldquo;</span><span
            class="c14">DLT</span><span class="c4">&rdquo;).</span></p>
          <p class="c1"><span
            class="c4">In order to use the Services you are required to register with MetaMask (for the login process), use the browser Chrome and have Libre Office as well as the Libre Office Plugin installed.</span>
          </p>
          <p class="c1"><span
            class="c4">Proxeus hosts its Services with an external service provider, currently BlockFactory AG, in Cham, Switzerland (</span><span
            class="c17 c4"><a class="c23"
                              target="_blank" href="https://www.blockfactory.com">Link</a></span><span
            class="c4">). Proxeus shall reserve the right to change its service provider at its sole discretion.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="2">
            <li class="c11"><h2 style="display:inline"><span class="c2">Service Content</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">Content created with the Services including the User Content as defined below (hereinafter &ldquo;</span><span
            class="c14">Service Content</span><span
            class="c4">&rdquo;) will be stored within the system of the service provider. Although Proxeus takes the best possible care in selecting the service provider, Proxeus shall not be held liable for damages or loss of data.</span>
          </p>
          <p class="c1"><span
            class="c4">Service Content shall only be accessible on the Proxeus platform, the User has no claim to an export function or any other way of downloading their Service Content. If Proxeus decides to make available such a function for certain Services, it is at Proxeus sole discretion and does not imply a right of availability for other Services.</span>
          </p>
          <p class="c1"><span
            class="c4">Some of the Services allow the integration of third party content (e.g., through template uploads or data entry in form</span><span
            class="c4">s</span><span
            class="c4">). Proxeus has no control in regards to completeness, validity, legality, correctness, quality and suitability for particular purposes of such third party content. Proxeus therefore shall not be responsible for service errors resulting from such third party content or any damage connected to such third party content.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="3">
            <li class="c11"><h2 style="display:inline"><span class="c2">Service security</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">You acknowledge that the Service Content is stored without any encryption or other security measures by Proxeus with the service provider. For the security measures provided by the service provider please consult the terms of service and the privacy statement of the service provider.</span>
          </p>
          <p class="c1"><span
            class="c4">By using the Service provided by Proxeus to make Service Content available to other Users, the User acknowledges that such sharing occurs via a deep link. By nature of a deep link every User in possession of such deep link will be able to use the Service Content. Proxeus shall not protect such links in any way.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="4">
            <li class="c11"><h2 style="display:inline"><span class="c2">Demo version</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">You are aware of the Services provided by Proxeus being a demo version (the &ldquo;</span><span
            class="c14">Demo</span><span
            class="c4">&rdquo;). Therefore, using the Services may require additional efforts and care on your side, while a continuing high quality of service delivery may not be guaranteed.</span>
          </p>
          <p class="c1"><span
            class="c4">Furthermore Proxeus shall not be responsible for any harm done to your hard- / software, privacy of your Ethereum identity or any other identity used while using the Demo.</span>
          </p>
          <p class="c1"><span
            class="c4">Proxeus also points out that Service Content created on the Demo may be subject to resets of the entire platform and may then not be available anymore. We strongly recommend not to use the Demo to build productive content on which you depend on. Proxeus shall not be liable for any damage resulting in loss of Service Content or any data after completion of the Demo.</span>
          </p>
          <p class="c1"><span
            class="c4">During the availability of the Demo, the Services are provided only in a hosted environment and interact solely with the Ropsten Test Net (Ethereum Testnet). Proxeus shall reserve the right to switch to the Ethereum main net or any another Ethereum network at their sole discretion.</span>
          </p>
          <p class="c1"><span
            class="c4">For the use of the Demo you may be provided with XES tokens which are reliant on the Ropsten Test Net, to see how the ecosystem manages payments – the first service being the registration of a document to the Ethereum environment for proof of existence and genuineness. Those XES tokens are only available in the Demo version (and on the Ropsten Test Net) and do not hold any value or right of usage of the Services on other Proxeus services environments. Due to the nature of the Ropsten Test Net, loss of XES tokens and the attached Service Content is possible and probable. Proxeus shall not be held liable for such a loss.</span><span
            class="c4">.</span></p>
          <ol class="c8 lst-kix_list_51-1" start="5">
            <li class="c11"><h2 style="display:inline"><span class="c2">Service disruption</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">Proxeus aims to ensure the availability of the Services. However the User acknowledges, that for technical reasons as well as due to the service provider&#39;s dependence on external factors, (e.g. unavailability of telecommunications networks, electricity outages, hardware and/or software failure etc.), the uninterrupted availability of the Services cannot be guaranteed. The User shall therefore not assert a claim for continual access to the Service. Access restrictions of any nature shall not constitute grounds for warranty claims.</span>
          </p>
          <ol class="c8 lst-kix_list_51-0" start="4">
            <li class="c7"><h1 style="display:inline"><span>Privacy policy and data protection</span></h1></li>
          </ol>
          <p class="c1"><span
            class="c4">Besides the Service Content (as described in section 3.1) we collect some personal data / User data to create and administer your account as well as data collected during your usage of the Services.</span>
          </p>
          <p class="c1"><span
            class="c4">You agree and confirm that you consent that Proxeus can process your personal data / User data for its own purposes such as to further develop the Services, administrate our platform and comply with all legal requirements. Proxeus ensures that the personal data of its Users is only collected, stored and processed to the extent required to render and further develop the Services and to the extent permitted by the applicable data protection laws. Further information on data processing and data protection is provided in the Proxeus </span><span
            class="c4">Privacy Statement</span><span class="c4">&nbsp;(</span><span
            class="c4">available</span><span class="c4">&nbsp;</span><span class="c4 c17"><a class="c23"
                                                                                             target="_blank"
                                                                                             href="https://proxeus.com/en/privacy-policy/">here</a></span><span
            class="c4">).</span></p>
          <p class="c1"><span
            class="c4">The User confirms that all the personal data provided by him or her is true and complete. Proxeus reserves the right to request appropriate proof of identity on a case by case basis.</span>
          </p>
          <ol class="c8 lst-kix_list_51-0" start="5">
            <li class="c10"><h1 style="display:inline"><span
              class="c6 c12">Representations, warranties and risks</span></h1></li>
          </ol>
          <ol class="c8 lst-kix_list_51-1 start" start="1">
            <li class="c11"><h2 style="display:inline"><span class="c2">Warranty disclaimer</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">You expressly understand and agree that your use of the Services is at your sole risk. The Services are provided on an &ldquo;AS IS&rdquo; and &ldquo;as available&rdquo; basis, without warranties of any kind, either express or implied, including, without limitation, implied warranties of merchantability, non-infringement or fitness for a particular purpose. You acknowledge that Proxeus has no control over, and no duty to take any action regarding the following: which Users gain access to or use the Services; what effects the Service Content may have on you; how you may interpret or use the Service Content; or what actions you may take as a result of having been exposed to the Service Content. Proxeus makes no representations concerning any Service Content contained in or accessed through the Services, and Proxeus shall not be responsible or liable for the accuracy, copyright compliance, legality or decency of material contained in or accessed through the Services.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="2">
            <li class="c11"><h2 style="display:inline"><span class="c2">Risk of cryptographic systems</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">By interacting with the Services in any way, you represent that you understand the inherent risks associated with cryptographic systems; and warrant that you have an understanding of the usage and intricacies of native cryptographic tokens, like Ether (ETH), smart contract based tokens such as those that follow the Ethereum Token Standard, and blockchain based software systems.</span>
          </p>
          <p class="c1" id="h.30j0zll"><span
            class="c4">You acknowledge and understand that cryptography is a progressing field. Advances in code cracking or technical advances such as the development of quantum computers may present risks to cryptocurrencies, which could result in the theft or loss of your cryptogr</span><span
            class="c4">aphic tokens or property. Proxeus relies fully on state-of-the-art third party cryptographic solutions and will, to the extent possible, update these third party solutions underlying the Services in a timely manner. Proxeus does not</span><span
            class="c4">&nbsp;guarantee or otherwise represent full security of the system. By every use of the Services or at each time you access Service Content, you acknowledge your awareness of these inherent risks.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="3">
            <li class="c11"><h2 style="display:inline"><span
              class="c2">Risk of regulatory actions in one or more jurisdiction</span>
            </h2></li>
          </ol>
          <p class="c1"><span
            class="c4">You acknowledge and understand that Proxeus or the underling DLT could be impacted by one or more regulatory inquiries or regulatory action, which could impede or limit the ability of Proxeus to continue to develop, or which could impede or limit your ability to access or use the Services or the DLT. By every use of the Services or at each time you access Service Content, you acknowledge your awareness of these inherent risks.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="4">
            <li class="c11"><h2 style="display:inline"><span class="c2">Volatility of cryptocurrencies</span></h2>
            </li>
          </ol>
          <p class="c1"><span
            class="c4">You understand that the monetary value of Ethereum and other DLT&rsquo;s and associated currencies or tokens are highly volatile due to many factors including but not limited to technology, speculation and security risks. You also acknowledge that the cost of transacting on such technologies is variable and may increase at any time causing impact to any activities taking place on the DLT. You acknowledge these risks and agree that Proxeus shall not be liable for such fluctuations or increased costs.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="5">
            <li class="c11"><h2 style="display:inline"><span class="c2">Application security</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">You acknowledge that such applications are code subject to flaws and acknowledge that you are solely responsible for evaluating any code provided by the Services or Service Content and the trustworthiness of any third-party websites, products, smart-contracts, or Service Content you access or use through the Services. You further expressly acknowledge and represent that such applications can be written maliciously or negligently, that Proxeus shall not be liable for your interaction with such applications and that such applications may cause the loss of property or even identity. This warning and others later provided by Proxeus are in no way evidence or represent an on-going duty to alert you to all of the potential risks of utilizing the Services or Service Content.</span>
          </p>
          <ol class="c8 lst-kix_list_51-0" start="6">
            <li class="c7"><h1 style="display:inline"><span class="c6 c12">General user obligations</span></h1></li>
          </ol>
          <ol class="c8 lst-kix_list_51-1 start" start="1">
            <li class="c11"><h2 style="display:inline"><span class="c2">General</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">You are solely responsible for complying with all laws that apply to you and to your end-users. This also applies explicitly to any additional legal regulations regarding the operation of a service in connection with cryptocurrencies.</span>
          </p>
          <p class="c1"><span
            class="c4">You are solely responsible for saving and securing your Service Content.</span></p>
          <p class="c1"><span
            class="c4">In the event of claims (e.g. secondary liability, third party liability etc.) due to illegal content which the User has placed on the Proxeus platform, the User undertakes to indemnify Proxeus and hold Proxeus harmless from any claim or demand (including and not limited to cease and desist orders with contractual fines, revocations, damages, rectifications etc.). In such a case, the User shall be under the obligation to assist Proxeus in every manner in responding to and in the defence of such claims.</span>
          </p>
          <ol class="c8 lst-kix_list_51-1" start="2">
            <li class="c11"><h2 style="display:inline"><span class="c2">Prohibited acts</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">The User is prohibited from undertaking any acts on the Proxeus platform or by using the Services (e.g. submit, transmit or display any User content) that violate the law, and/or the rights of third parties or that violate the basic principles regarding the protection of minors.</span>
          </p>
          <p class="c1"><span
            class="c4">The publication, spreading, offering and the advertisement or advertising of the following content in particular is expressly prohibited:</span>
          </p>
          <ul class="c8 lst-kix_ek190fjkpfvz-0 start">
            <li class="c20"><span class="c4">content that is pornographic, obscene or immoral in nature;</span></li>
            <li class="c3"><span
              class="c4">content that is in violation of relevant legislation regarding the protection of minors, that violate data protection regulations and that otherwise violate the law or content/services/products that are fraudulent in nature;</span>
            </li>
            <li class="c3"><span
              class="c4">content that glorifies or trivialises war, terror and other acts of violence against people or animals;</span>
            </li>
            <li class="c3"><span
              class="c4">content on sex, race, colour, ethnic group or social origin, language, religion or belief, political or any other opinion, birth, disability, age or sexual orientation, nationality etc. that may insult or slander other Users or third parties;</span>
            </li>
            <li class="c3"><span
              class="c4">content that is deemed to promote or support racism, radicalism, fascism, fanaticism, hate, physical and psychological violence or illegal activity (whether explicit or implicit) or that otherwise breach the standards of common decency;</span>
            </li>
            <li class="c3"><span
              class="c4">content directed at insulting or slandering or defaming (defamation) other participants / persons or third parties;</span>
            </li>
            <li class="c13"><span
              class="c4">content, services and or products that are legally protected or encumbered with the rights of third parties (e.g. copyright/trademark protection), without being demonstrably entitled to do so.</span>
            </li>
          </ul>
          <p class="c1"><span
            class="c4">This obligation also applies to any links (Hyperlinks) included by the User, pertaining to such content as described above, included on external platforms or services.</span>
          </p>
          <p class="c1"><span
            class="c4">Furthermore, independent of any possible legal ramifications, the following activities as concerns the setting of the User&rsquo;s own content (e.g. via the setting of links), are also prohibited:</span>
          </p>
          <ul class="c8 lst-kix_ek190fjkpfvz-0">
            <li class="c3"><span
              class="c6 c4">distribution or other transmission or execution of viruses, trojans and other damaging data;</span>
            </li>
            <li class="c3"><span
              class="c6 c4">sending of junk, spam or scam mails as well as &ldquo;chainmail&rdquo;;</span></li>
            <li class="c3"><span class="c4 c6">utilisation of the Services for spam purposes or SEO spam;</span></li>
            <li class="c3"><span
              class="c6 c4">harassment of others e.g. via multiple instances of contact via the Services without or in contravention to the reaction of the recipient, as well as the promotion or encouragement of such harassment;</span>
            </li>
            <li class="c3"><span
              class="c6 c4">requesting of passwords or other personal data for commercial or illegal purposes (phishing);</span>
            </li>
            <li class="c3"><span
              class="c6 c4">making available, publication, offering or advertisement of bonus systems, Paid4Mail-Services, PopUp-Services, Snowball systems, pyramid schemes or other similarly functioning schemes as well as engaging in unethical marketing or advertising;</span>
            </li>
            <li class="c3"><span
              class="c6 c4">offering of games of chance such as public sports betting, general betting and lotteries etc. without permission from the relevant authority;</span>
            </li>
            <li class="c3"><span
              class="c6 c4">procurement or indirect procurement of loans and private loans in a commercial manner without demonstrable permission from a relevant authority.</span>
            </li>
          </ul>
          <p class="c1"><span
            class="c4">Furthermore, every action that could influence the regular functioning of the Proxeus platform or the Services is prohibited.</span>
          </p>
          <ol class="c8 lst-kix_list_54-0 start" start="7">
            <li class="c10"><h1 style="display:inline"><span class="c6 c12">Intellectual property rights</span></h1>
            </li>
          </ol>
          <ol class="c8 lst-kix_list_54-1 start" start="1">
            <li class="c11"><h2 style="display:inline"><span class="c2">Responsibility for User Content</span></h2>
            </li>
          </ol>
          <p class="c1"><span
            class="c4">We offer Users the possibility to upload content and use this content in the Services as well as to make said content available to others (hereinafter &ldquo;</span><span
            class="c14">User Content</span><span class="c4">&rdquo;).</span><span class="c4">&nbsp;</span><span
            class="c4">In order for Proxeus to perform the Services the User grants free usage rights to Proxeus for the publication, editing and public broadcasting of the User Content. This does not entitle Proxeus to sell the User Content to a third party. The copyright of the author shall remain unaffected.</span>
          </p>
          <p class="c1"><span
            class="c4">With the uploading of such content the User grants Proxeus gratuitous usage rights to the various content, including:</span>
          </p>
          <ul class="c8 lst-kix_ek190fjkpfvz-0">
            <li class="c3"><span
              class="c6 c4">to store the User Content on the Proxeus servers as well as service provider servers or to publish and make the User Content publicly available.</span>
            </li>
            <li class="c3"><span
              class="c4">to edit and reproduce related data to the extent necessary for delivery and publication of the User Content.</span>
            </li>
          </ul>
          <p class="c1"><span
            class="c4">The Users are solely responsible for all content they upload. Proxeus is not obliged to inspect the User Content as regards of completeness, validity, legality, correctness, quality and suitability for particular </span><span
            class="c4">purposes</span><span class="c4">.</span></p>
          <p class="c1"><span
            class="c4">Proxeus reserves the right to refuse, block or remove User Content without prior notice, should the User Content constitute or lead to a violation of these Terms, or where there is clear evidence that a serious violation of these Terms shall result. The legitimate interests of the User shall be taken into account and appropriate measures to prevent and / or remove the violation shall be taken.</span>
          </p>
          <ol class="c8 lst-kix_list_54-1" start="2">
            <li class="c11"><h2 style="display:inline"><span class="c2">Sharing of Service Content</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">By using the share functionality provided by Proxeus to make the Service Content available to other Users, the User grants gratuitous usage rights to the receiving User for the publication, editing and public broadcasting of the Service </span><span
            class="c4">Content</span><span class="c4">.</span></p>
          <ol class="c8 lst-kix_list_54-0" start="8">
            <li class="c7"><h1 style="display:inline"><span class="c6 c12">Limitations of liability</span></h1></li>
          </ol>
          <p class="c1"><span
            class="c4">Without prejudice to our obligations under these Terms, we do not warrant that the Services will meet any particular requirements or that their operation will be entirely error-free or that defects are capable of correction or improvement.</span>
          </p>
          <p class="c1"><span
            class="c4">Proxeus&rsquo;s liability towards the Users is excluded to the maximum extent permitted by the applicable law. In particular, neither Proxeus nor its affiliates, agents or subcontractors shall be liable to the User or any third party for the following loss or damage, whether arising in tort, contract, breach of statutory duty or otherwise, and even if foreseeable by Proxeus: any indirect, special, consequential or incidental loss of profits, business, contracts, goodwill, reputation, opportunity, revenue production, or anticipated savings howsoever caused, arising out of, or in connection with, any supply, failure to supply or delay in supplying any of the Services or otherwise in connection with these Terms (including fundamental breach or breach of a fundamental term) or any other theory of law.</span>
          </p>
          <ol class="c8 lst-kix_list_54-0" start="9">
            <li class="c10"><h1 style="display:inline"><span class="c6 c12">Miscellaneous provisions</span></h1></li>
          </ol>
          <ol class="c8 lst-kix_list_54-1 start" start="1">
            <li class="c11"><h2 style="display:inline"><span class="c2">Severability clause</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">Should any part or provision of these Terms be held to be invalid or unenforceable by any competent arbitral tribunal, court, governmental or administrative authority having jurisdiction, the other provisions of the Terms shall nonetheless remain valid.</span>
          </p>
          <ol class="c8 lst-kix_list_54-1" start="2">
            <li class="c11"><h2 style="display:inline"><span class="c2">Amendments</span></h2></li>
          </ol>
          <p class="c1"><span
            class="c4">Proxeus retains the right to amend these Terms at any time, also within the current contractual relationships. Proxeus will make such changes public by updating these Terms. All changes come into effect immediately after publication. Your continued use of our Services after the publication of the new Terms means that you accept and agree to the new Terms. You are therefore advised to check the Terms </span><span
            class="c4">regularly</span><span class="c4">.</span></p>
          <ol class="c8 lst-kix_list_54-1" start="3">
            <li class="c11"><h2 style="display:inline"><span>Jurisdiction and applicable </span><span>law</span></h2>
            </li>
          </ol>
          <p class="c22"><span class="c4">Subject to mandatory statutory law, t</span><span
            class="c6 c4">hese Terms shall be governed by and construed in accordance with the substantive laws of</span><span
            class="c4">&nbsp;the Principality of Liechtenstein</span><span class="c6 c4">.</span></p>
          <p class="c22" id="h.1fob9te"><span class="c4">Subject to mandatory statutory law, a</span><span
            class="c6 c4">ny disputes arising out of or in connection with these Terms, including any disputes regarding their validity or termination and the validity of this choice of forum provision shall exclusively be brought to the ordinary courts of the Principality o</span><span
            class="c4">f Liechtenstein</span><span class="c6 c4">, venue being </span><span
            class="c4">Vaduz</span><span class="c6 c4">, Principality of </span><span
            class="c4">Liechtenstein</span><span class="c6 c4">.</span></p>
        </div>
        <div class="modal-footer">
          <button type="button" @click="acceptTC" class="btn btn-primary">I agree</button>
          <button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import Spinner from '../components/Spinner'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'login',
  components: {
    Spinner
  },
  data () {
    return {
      account: undefined,
      pwlogin: true,
      email: '',
      password: '',
      hasError: false,
      walletErrorMessage: '',
      loginErrorMessage: '',
      metamaskLoginAvailable: false,
      signing: false
    }
  },
  created () {
    if (window.web3 !== undefined) {
      this.metamaskLoginAvailable = true
    }
  },
  mounted () {
    this.$refs.inputEmail && this.$refs.inputEmail.focus()
  },
  computed: {
    wallet () {
      return this.app.wallet
    }
  },
  methods: {
    login () {
      this.account = this.email
      this.pwlogin = true
      if (this.checkTermsAndConditions()) {
        axios.post('/api/login', { email: this.email, password: this.password }).then(res => {
          this.hasError = false
          window.location = res.data.location || '/admin/workflow'
        }, (err) => {
          this.loginErrorMessage = 'You have entered an invalid username or password'
          this.hasError = true
          this.$nextTick(() => {
            this.$refs.inputEmail.focus()
          })
          this.app.handleError(err)
        })
      }
    },
    checkTermsAndConditions () {
      let rememberAccept = localStorage.getItem('acc_' + this.account)
      if (rememberAccept && rememberAccept === 'yes') {
        return true
      }
      $(this.$refs.tcModal).modal('show')
      $(this.$refs.tcModal).on('hide.bs.modal', (e) => {
        this.declineTC()
      })
    },
    declineTC () {
      this.signing = false
    },
    acceptTC () {
      localStorage.setItem('acc_' + this.account, 'yes')
      this.app.acknowledgeFirstLogin()
      $(this.$refs.tcModal).modal('hide')
      if (this.pwlogin) {
        this.login()
      } else {
        this.metamaskLogin()
      }
    },
    async metamaskLogin () {
      if (window.ethereum) {
        try {
          await window.ethereum.enable()
          await this.wallet.wallet.setupDefaultAccount()
        } catch (e) {
          this.walletErrorMessage = 'Please grant access to MetaMask.'
          return
        }
      }
      axios.get('/api/challenge').then((response) => {
        this.signing = true
        let challenge = response.data
        let account = this.wallet.getCurrentAddress()
        this.account = account
        if (this.account === undefined) {
          this.walletErrorMessage = 'Please log in to MetaMask.'
          this.signing = false
          return
        }
        this.pwlogin = false
        if (this.checkTermsAndConditions()) {
          console.log(this.account)
          this.wallet.signMessage(challenge, this.account).then((signature) => {
            axios.post('/api/login', { signature }).then((res) => {
              if (res.status >= 200 && res.status <= 299) {
                window.location = res.data.location || '/admin/workflow'
              } else {
                this.walletErrorMessage = 'Could not verify signature.'
              }
            }, (err) => {
              this.signing = false
              this.walletErrorMessage = 'Could not verify signature.'
            })
          }).catch(() => {
            this.signing = false
            this.walletErrorMessage = 'Could not Sign Message.'
          })
        }
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: 'Error',
          text: 'Login failed. Please try again or if the error persists contact the platform operator.',
          type: 'warning'
        })
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "../assets/styles/variables.scss";

  @media print {
    * {
      display: block !important;
      visibility: visible !important;
      opacity: 1 !important;
      height: 100% !important;
      min-height: 100% !important;
      overflow: visible !important;
    }
    .modal-body {
      max-height: 100%;
      overflow: visible;
    }

    .modal-backdrop {
      display: none;
    }

    .modal {
      position: relative;
      display: block;
      opacity: 1;
      overflow: visible;
    }

    .modal-open {
      overflow: visible;
      overflow: visible;
    }

    .modal-dialog {
      transform: none;
      overflow: visible;
    }
  }

  .modal-body {
    h1 {
      font-size: $h1-font-size;
    }

    h2 {
      font-size: $h2-font-size;
    }

    h3 {
      font-size: $h3-font-size;
    }

    ol {
      padding-left: 0;
      margin-bottom: 0;
    }
  }

  h2 {
    padding-top: 1rem;
    font-size: 1.5rem !important;
  }

  h3 {
    padding-top: .5rem;
    font-size: 1.2rem !important;
  }

  .c10:before, .c7:before {
    font-size: $h1-font-size;
  }

  .c11:before {
    font-size: $h2-font-size;
  }

  .modal-header {
    background: $light;
  }

  .modal-footer {
    background: $light;
    border-bottom-right-radius: $border-radius;
    border-bottom-left-radius: $border-radius;
  }

  .modal-body {
    overflow-y: auto;
    max-height: calc(100vh - 280px);
  }

  .login-form {
    overflow: auto;
    margin: 0 auto;
    margin-top: 50px;
    height: 100%;
    max-width: 600px;
    padding-top: 40px;
    padding-bottom: 40px;
    border-radius: $border-radius;
  }

  .login-form-sm {
    max-width: 350px;
  }

  .form-signin {
    width: 100%;
    max-width: 330px;
    padding: 2rem;
    margin: 0 auto;
    z-index: 1000;

    .checkbox {
      font-weight: 400;
    }

    .form-control {
      position: relative;
      box-sizing: border-box;
      height: auto;
      padding: 10px;
      font-size: 16px;
    }
  }

  .tc-list {
    margin-bottom: 0;
    margin-top: 1.5rem;
    li {
      padding-left: .5rem;
      padding-bottom: .75rem;
    }
  }

  .modal-body {

    ol, li {
      list-style-type: none;
    }

    .lst-kix_list_33-5 > li {
      counter-increment: lst-ctn-kix_list_33-5
    }

    .lst-kix_list_21-8 > li {
      counter-increment: lst-ctn-kix_list_21-8
    }

    .lst-kix_list_5-0 > li {
      counter-increment: lst-ctn-kix_list_5-0
    }

    ol.lst-kix_list_2-3.start {
      counter-reset: lst-ctn-kix_list_2-3 0
    }

    ol.lst-kix_list_37-2.start {
      counter-reset: lst-ctn-kix_list_37-2 0
    }

    .lst-kix_list_45-2 > li {
      counter-increment: lst-ctn-kix_list_45-2
    }

    .lst-kix_list_42-1 > li:before {
      content: "" counter(lst-ctn-kix_list_42-1, decimal) " "
    }

    .lst-kix_list_40-1 > li {
      counter-increment: lst-ctn-kix_list_40-1
    }

    ol.lst-kix_list_5-3.start {
      counter-reset: lst-ctn-kix_list_5-3 0
    }

    .lst-kix_list_38-6 > li {
      counter-increment: lst-ctn-kix_list_38-6
    }

    .lst-kix_list_4-3 > li {
      counter-increment: lst-ctn-kix_list_4-3
    }

    .lst-kix_list_42-7 > li:before {
      content: "" counter(lst-ctn-kix_list_42-7, decimal) " "
    }

    .lst-kix_list_51-1 > li {
      counter-increment: lst-ctn-kix_list_51-1
    }

    ol.lst-kix_list_17-1.start {
      counter-reset: lst-ctn-kix_list_17-1 0
    }

    .lst-kix_list_42-5 > li:before {
      content: "" counter(lst-ctn-kix_list_42-5, decimal) " "
    }

    .lst-kix_list_42-3 > li:before {
      content: "" counter(lst-ctn-kix_list_42-3, decimal) " "
    }

    .lst-kix_list_32-8 > li {
      counter-increment: lst-ctn-kix_list_32-8
    }

    .lst-kix_list_24-8 > li:before {
      content: "" counter(lst-ctn-kix_list_24-8, lower-roman) ". "
    }

    .lst-kix_list_24-2 > li:before {
      content: "(" counter(lst-ctn-kix_list_24-2, lower-latin) ") "
    }

    ol.lst-kix_list_54-3.start {
      counter-reset: lst-ctn-kix_list_54-3 0
    }

    .lst-kix_list_24-4 > li:before {
      content: "" counter(lst-ctn-kix_list_24-4, decimal) ") "
    }

    ol.lst-kix_list_34-2.start {
      counter-reset: lst-ctn-kix_list_34-2 0
    }

    .lst-kix_list_9-4 > li {
      counter-increment: lst-ctn-kix_list_9-4
    }

    .lst-kix_list_24-6 > li:before {
      content: "" counter(lst-ctn-kix_list_24-6, lower-roman) ") "
    }

    ol.lst-kix_list_14-1.start {
      counter-reset: lst-ctn-kix_list_14-1 0
    }

    .lst-kix_list_23-6 > li:before {
      content: "" counter(lst-ctn-kix_list_23-6, decimal) ". "
    }

    .lst-kix_list_23-2 > li:before {
      content: "" counter(lst-ctn-kix_list_23-2, lower-roman) ") "
    }

    .lst-kix_list_23-0 > li:before {
      content: "" counter(lst-ctn-kix_list_23-0, decimal) ") "
    }

    .lst-kix_list_23-8 > li:before {
      content: "" counter(lst-ctn-kix_list_23-8, lower-roman) ". "
    }

    .lst-kix_list_3-6 > li {
      counter-increment: lst-ctn-kix_list_3-6
    }

    ol.lst-kix_list_44-0.start {
      counter-reset: lst-ctn-kix_list_44-0 0
    }

    .lst-kix_list_24-0 > li:before {
      content: "Article " counter(lst-ctn-kix_list_24-0, upper-roman) ". "
    }

    .lst-kix_list_23-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_23-4, lower-latin) ") "
    }

    ol.lst-kix_list_51-3.start {
      counter-reset: lst-ctn-kix_list_51-3 0
    }

    .lst-kix_list_43-3 > li:before {
      content: "" counter(lst-ctn-kix_list_43-3, decimal) " "
    }

    ol.lst-kix_list_31-2.start {
      counter-reset: lst-ctn-kix_list_31-2 0
    }

    .lst-kix_list_22-2 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) "." counter(lst-ctn-kix_list_22-2, decimal) ". "
    }

    .lst-kix_list_22-6 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) "." counter(lst-ctn-kix_list_22-2, decimal) "." counter(lst-ctn-kix_list_22-3, decimal) "." counter(lst-ctn-kix_list_22-4, decimal) "." counter(lst-ctn-kix_list_22-5, decimal) "." counter(lst-ctn-kix_list_22-6, decimal) ". "
    }

    .lst-kix_list_43-7 > li:before {
      content: "" counter(lst-ctn-kix_list_43-7, decimal) " "
    }

    .lst-kix_list_43-1 > li:before {
      content: "" counter(lst-ctn-kix_list_43-1, decimal) " "
    }

    .lst-kix_list_22-0 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) ". "
    }

    .lst-kix_list_22-8 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) "." counter(lst-ctn-kix_list_22-2, decimal) "." counter(lst-ctn-kix_list_22-3, decimal) "." counter(lst-ctn-kix_list_22-4, decimal) "." counter(lst-ctn-kix_list_22-5, decimal) "." counter(lst-ctn-kix_list_22-6, decimal) "." counter(lst-ctn-kix_list_22-7, decimal) "." counter(lst-ctn-kix_list_22-8, decimal) ". "
    }

    ol.lst-kix_list_22-5.start {
      counter-reset: lst-ctn-kix_list_22-5 0
    }

    .lst-kix_list_5-7 > li {
      counter-increment: lst-ctn-kix_list_5-7
    }

    ol.lst-kix_list_34-7.start {
      counter-reset: lst-ctn-kix_list_34-7 0
    }

    .lst-kix_list_43-5 > li:before {
      content: "" counter(lst-ctn-kix_list_43-5, decimal) " "
    }

    .lst-kix_list_22-4 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) "." counter(lst-ctn-kix_list_22-2, decimal) "." counter(lst-ctn-kix_list_22-3, decimal) "." counter(lst-ctn-kix_list_22-4, decimal) ". "
    }

    ol.lst-kix_list_32-8.start {
      counter-reset: lst-ctn-kix_list_32-8 0
    }

    ol.lst-kix_list_25-5.start {
      counter-reset: lst-ctn-kix_list_25-5 0
    }

    .lst-kix_list_41-7 > li:before {
      content: "" counter(lst-ctn-kix_list_41-7, decimal) " "
    }

    ol.lst-kix_list_24-4.start {
      counter-reset: lst-ctn-kix_list_24-4 0
    }

    .lst-kix_list_6-4 > li {
      counter-increment: lst-ctn-kix_list_6-4
    }

    .lst-kix_list_41-1 > li:before {
      content: "" counter(lst-ctn-kix_list_41-1, decimal) " "
    }

    .lst-kix_list_40-7 > li:before {
      content: "" counter(lst-ctn-kix_list_40-7, decimal) " "
    }

    ol.lst-kix_list_27-4.start {
      counter-reset: lst-ctn-kix_list_27-4 0
    }

    .lst-kix_list_40-5 > li:before {
      content: "" counter(lst-ctn-kix_list_40-5, decimal) " "
    }

    ol.lst-kix_list_15-2.start {
      counter-reset: lst-ctn-kix_list_15-2 0
    }

    ol.lst-kix_list_42-1.start {
      counter-reset: lst-ctn-kix_list_42-1 0
    }

    .lst-kix_list_41-5 > li:before {
      content: "" counter(lst-ctn-kix_list_41-5, decimal) " "
    }

    ol.lst-kix_list_37-7.start {
      counter-reset: lst-ctn-kix_list_37-7 0
    }

    ol.lst-kix_list_2-8.start {
      counter-reset: lst-ctn-kix_list_2-8 0
    }

    .lst-kix_list_41-3 > li:before {
      content: "" counter(lst-ctn-kix_list_41-3, decimal) " "
    }

    ol.lst-kix_list_5-8.start {
      counter-reset: lst-ctn-kix_list_5-8 0
    }

    .lst-kix_list_1-3 > li {
      counter-increment: lst-ctn-kix_list_1-3
    }

    .lst-kix_list_40-3 > li:before {
      content: "" counter(lst-ctn-kix_list_40-3, decimal) " "
    }

    ol.lst-kix_list_12-2.start {
      counter-reset: lst-ctn-kix_list_12-2 0
    }

    .lst-kix_list_40-1 > li:before {
      content: "" counter(lst-ctn-kix_list_40-1, decimal) " "
    }

    .lst-kix_list_42-2 > li {
      counter-increment: lst-ctn-kix_list_42-2
    }

    ol.lst-kix_list_38-3.start {
      counter-reset: lst-ctn-kix_list_38-3 0
    }

    .lst-kix_list_24-8 > li {
      counter-increment: lst-ctn-kix_list_24-8
    }

    ol.lst-kix_list_3-4.start {
      counter-reset: lst-ctn-kix_list_3-4 0
    }

    .lst-kix_list_21-8 > li:before {
      content: "" counter(lst-ctn-kix_list_21-8, decimal) ". "
    }

    ol.lst-kix_list_18-2.start {
      counter-reset: lst-ctn-kix_list_18-2 0
    }

    .lst-kix_list_26-8 > li:before {
      content: "" counter(lst-ctn-kix_list_26-8, lower-roman) ". "
    }

    .lst-kix_list_21-0 > li:before {
      content: "" counter(lst-ctn-kix_list_21-0, decimal) ". "
    }

    .lst-kix_list_13-1 > li {
      counter-increment: lst-ctn-kix_list_13-1
    }

    .lst-kix_list_26-4 > li:before {
      content: "" counter(lst-ctn-kix_list_26-4, decimal) ") "
    }

    ol.lst-kix_list_36-1.start {
      counter-reset: lst-ctn-kix_list_36-1 0
    }

    .lst-kix_list_42-4 > li {
      counter-increment: lst-ctn-kix_list_42-4
    }

    .lst-kix_list_21-4 > li:before {
      content: "" counter(lst-ctn-kix_list_21-4, decimal) ". "
    }

    .lst-kix_list_26-0 > li:before {
      content: "Article " counter(lst-ctn-kix_list_26-0, upper-roman) ". "
    }

    .lst-kix_list_31-4 > li {
      counter-increment: lst-ctn-kix_list_31-4
    }

    .lst-kix_list_31-2 > li {
      counter-increment: lst-ctn-kix_list_31-2
    }

    ol.lst-kix_list_38-8.start {
      counter-reset: lst-ctn-kix_list_38-8 0
    }

    ol.lst-kix_list_19-5.start {
      counter-reset: lst-ctn-kix_list_19-5 0
    }

    .lst-kix_list_35-8 > li {
      counter-increment: lst-ctn-kix_list_35-8
    }

    ol.lst-kix_list_39-6.start {
      counter-reset: lst-ctn-kix_list_39-6 0
    }

    ol.lst-kix_list_45-6.start {
      counter-reset: lst-ctn-kix_list_45-6 0
    }

    ol.lst-kix_list_26-3.start {
      counter-reset: lst-ctn-kix_list_26-3 0
    }

    .lst-kix_list_45-5 > li:before {
      content: "" counter(lst-ctn-kix_list_45-5, decimal) " "
    }

    .lst-kix_list_25-0 > li:before {
      content: "" counter(lst-ctn-kix_list_25-0, decimal) ") "
    }

    ol.lst-kix_list_10-8.start {
      counter-reset: lst-ctn-kix_list_10-8 0
    }

    .lst-kix_list_39-0 > li:before {
      content: "" counter(lst-ctn-kix_list_39-0, decimal) " "
    }

    ol.lst-kix_list_40-7.start {
      counter-reset: lst-ctn-kix_list_40-7 0
    }

    .lst-kix_list_37-0 > li {
      counter-increment: lst-ctn-kix_list_37-0
    }

    ol.lst-kix_list_21-4.start {
      counter-reset: lst-ctn-kix_list_21-4 0
    }

    .lst-kix_list_44-1 > li:before {
      content: "" counter(lst-ctn-kix_list_44-1, decimal) " "
    }

    .lst-kix_list_45-1 > li:before {
      content: "" counter(lst-ctn-kix_list_45-1, decimal) " "
    }

    ol.lst-kix_list_20-6.start {
      counter-reset: lst-ctn-kix_list_20-6 0
    }

    .lst-kix_list_13-8 > li {
      counter-increment: lst-ctn-kix_list_13-8
    }

    .lst-kix_list_2-2 > li {
      counter-increment: lst-ctn-kix_list_2-2
    }

    .lst-kix_list_44-5 > li:before {
      content: "" counter(lst-ctn-kix_list_44-5, decimal) " "
    }

    ol.lst-kix_list_4-7.start {
      counter-reset: lst-ctn-kix_list_4-7 0
    }

    .lst-kix_list_26-2 > li {
      counter-increment: lst-ctn-kix_list_26-2
    }

    .lst-kix_list_40-8 > li {
      counter-increment: lst-ctn-kix_list_40-8
    }

    .lst-kix_list_27-4 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) "." counter(lst-ctn-kix_list_27-2, decimal) "." counter(lst-ctn-kix_list_27-3, decimal) "." counter(lst-ctn-kix_list_27-4, decimal) ". "
    }

    .lst-kix_list_20-2 > li {
      counter-increment: lst-ctn-kix_list_20-2
    }

    .lst-kix_list_6-6 > li {
      counter-increment: lst-ctn-kix_list_6-6
    }

    ol.lst-kix_list_15-7.start {
      counter-reset: lst-ctn-kix_list_15-7 0
    }

    .lst-kix_list_13-6 > li {
      counter-increment: lst-ctn-kix_list_13-6
    }

    ol.lst-kix_list_14-6.start {
      counter-reset: lst-ctn-kix_list_14-6 0
    }

    .lst-kix_list_39-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_39-4, lower-roman) ") "
    }

    .lst-kix_list_39-8 > li:before {
      content: "" counter(lst-ctn-kix_list_39-8, decimal) " "
    }

    ol.lst-kix_list_26-8.start {
      counter-reset: lst-ctn-kix_list_26-8 0
    }

    .lst-kix_list_19-6 > li {
      counter-increment: lst-ctn-kix_list_19-6
    }

    .lst-kix_list_27-0 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) ". "
    }

    .lst-kix_list_20-4 > li {
      counter-increment: lst-ctn-kix_list_20-4
    }

    ol.lst-kix_list_44-5.start {
      counter-reset: lst-ctn-kix_list_44-5 0
    }

    ol.lst-kix_list_33-6.start {
      counter-reset: lst-ctn-kix_list_33-6 0
    }

    ol.lst-kix_list_20-1.start {
      counter-reset: lst-ctn-kix_list_20-1 0
    }

    .lst-kix_list_25-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_25-4, lower-latin) ") "
    }

    .lst-kix_list_19-4 > li {
      counter-increment: lst-ctn-kix_list_19-4
    }

    .lst-kix_list_35-1 > li {
      counter-increment: lst-ctn-kix_list_35-1
    }

    .lst-kix_list_46-5 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_24-1 > li {
      counter-increment: lst-ctn-kix_list_24-1
    }

    .lst-kix_list_25-8 > li:before {
      content: "" counter(lst-ctn-kix_list_25-8, lower-roman) ". "
    }

    ol.lst-kix_list_50-0.start {
      counter-reset: lst-ctn-kix_list_50-0 0
    }

    ol.lst-kix_list_51-8.start {
      counter-reset: lst-ctn-kix_list_51-8 0
    }

    .lst-kix_list_51-8 > li {
      counter-increment: lst-ctn-kix_list_51-8
    }

    ol.lst-kix_list_26-6.start {
      counter-reset: lst-ctn-kix_list_26-6 0
    }

    .lst-kix_list_44-5 > li {
      counter-increment: lst-ctn-kix_list_44-5
    }

    .lst-kix_list_37-2 > li {
      counter-increment: lst-ctn-kix_list_37-2
    }

    .lst-kix_list_15-2 > li {
      counter-increment: lst-ctn-kix_list_15-2
    }

    .lst-kix_list_22-5 > li {
      counter-increment: lst-ctn-kix_list_22-5
    }

    .lst-kix_list_28-8 > li:before {
      content: "" counter(lst-ctn-kix_list_28-8, lower-roman) ". "
    }

    .lst-kix_list_24-6 > li {
      counter-increment: lst-ctn-kix_list_24-6
    }

    .lst-kix_list_17-3 > li {
      counter-increment: lst-ctn-kix_list_17-3
    }

    .lst-kix_list_28-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_28-4, lower-latin) ") "
    }

    .lst-kix_list_39-3 > li {
      counter-increment: lst-ctn-kix_list_39-3
    }

    .lst-kix_list_28-3 > li {
      counter-increment: lst-ctn-kix_list_28-3
    }

    ol.lst-kix_list_25-0.start {
      counter-reset: lst-ctn-kix_list_25-0 0
    }

    .lst-kix_list_35-6 > li {
      counter-increment: lst-ctn-kix_list_35-6
    }

    ol.lst-kix_list_13-0.start {
      counter-reset: lst-ctn-kix_list_13-0 0
    }

    ol.lst-kix_list_32-3.start {
      counter-reset: lst-ctn-kix_list_32-3 0
    }

    ol.lst-kix_list_50-2.start {
      counter-reset: lst-ctn-kix_list_50-2 0
    }

    .lst-kix_list_11-5 > li {
      counter-increment: lst-ctn-kix_list_11-5
    }

    .lst-kix_list_28-0 > li:before {
      content: "" counter(lst-ctn-kix_list_28-0, decimal) ") "
    }

    .lst-kix_list_27-8 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) "." counter(lst-ctn-kix_list_27-2, decimal) "." counter(lst-ctn-kix_list_27-3, decimal) "." counter(lst-ctn-kix_list_27-4, decimal) "." counter(lst-ctn-kix_list_27-5, decimal) "." counter(lst-ctn-kix_list_27-6, decimal) "." counter(lst-ctn-kix_list_27-7, decimal) "." counter(lst-ctn-kix_list_27-8, decimal) ". "
    }

    ol.lst-kix_list_20-3.start {
      counter-reset: lst-ctn-kix_list_20-3 0
    }

    .lst-kix_list_4-1 > li {
      counter-increment: lst-ctn-kix_list_4-1
    }

    .lst-kix_list_19-1 > li:before {
      content: "" counter(lst-ctn-kix_list_19-1, decimal) ". "
    }

    .lst-kix_list_27-8 > li {
      counter-increment: lst-ctn-kix_list_27-8
    }

    .lst-kix_list_19-3 > li:before {
      content: "" counter(lst-ctn-kix_list_19-3, decimal) ". "
    }

    ol.lst-kix_list_38-0.start {
      counter-reset: lst-ctn-kix_list_38-0 0
    }

    .lst-kix_list_15-0 > li {
      counter-increment: lst-ctn-kix_list_15-0
    }

    ol.lst-kix_list_6-6.start {
      counter-reset: lst-ctn-kix_list_6-6 0
    }

    .lst-kix_list_47-5 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_39-5 > li {
      counter-increment: lst-ctn-kix_list_39-5
    }

    .lst-kix_list_32-6 > li {
      counter-increment: lst-ctn-kix_list_32-6
    }

    ol.lst-kix_list_29-6.start {
      counter-reset: lst-ctn-kix_list_29-6 0
    }

    .lst-kix_list_11-0 > li {
      counter-increment: lst-ctn-kix_list_11-0
    }

    ol.lst-kix_list_1-5.start {
      counter-reset: lst-ctn-kix_list_1-5 0
    }

    ol.lst-kix_list_9-6.start {
      counter-reset: lst-ctn-kix_list_9-6 0
    }

    ol.lst-kix_list_16-3.start {
      counter-reset: lst-ctn-kix_list_16-3 0
    }

    .lst-kix_list_22-7 > li {
      counter-increment: lst-ctn-kix_list_22-7
    }

    .lst-kix_list_47-3 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_4-5.start {
      counter-reset: lst-ctn-kix_list_4-5 0
    }

    .lst-kix_list_37-0 > li:before {
      content: "" counter(lst-ctn-kix_list_37-0, decimal) " "
    }

    .lst-kix_list_26-7 > li {
      counter-increment: lst-ctn-kix_list_26-7
    }

    .lst-kix_list_41-0 > li {
      counter-increment: lst-ctn-kix_list_41-0
    }

    .lst-kix_list_33-7 > li {
      counter-increment: lst-ctn-kix_list_33-7
    }

    .lst-kix_list_5-2 > li {
      counter-increment: lst-ctn-kix_list_5-2
    }

    .lst-kix_list_37-7 > li {
      counter-increment: lst-ctn-kix_list_37-7
    }

    .lst-kix_list_37-8 > li:before {
      content: "" counter(lst-ctn-kix_list_37-8, upper-latin) ". "
    }

    .lst-kix_list_37-6 > li:before {
      content: "" counter(lst-ctn-kix_list_37-6, upper-roman) ". "
    }

    .lst-kix_list_28-5 > li {
      counter-increment: lst-ctn-kix_list_28-5
    }

    ol.lst-kix_list_35-0.start {
      counter-reset: lst-ctn-kix_list_35-0 0
    }

    .lst-kix_list_22-0 > li {
      counter-increment: lst-ctn-kix_list_22-0
    }

    .lst-kix_list_46-1 > li:before {
      content: "o  "
    }

    ol.lst-kix_list_33-4.start {
      counter-reset: lst-ctn-kix_list_33-4 0
    }

    .lst-kix_list_21-6 > li {
      counter-increment: lst-ctn-kix_list_21-6
    }

    .lst-kix_list_10-3 > li {
      counter-increment: lst-ctn-kix_list_10-3
    }

    ol.lst-kix_list_1-0.start {
      counter-reset: lst-ctn-kix_list_1-0 0
    }

    ol.lst-kix_list_13-3.start {
      counter-reset: lst-ctn-kix_list_13-3 0
    }

    .lst-kix_list_51-6 > li {
      counter-increment: lst-ctn-kix_list_51-6
    }

    .lst-kix_list_26-0 > li {
      counter-increment: lst-ctn-kix_list_26-0
    }

    ol.lst-kix_list_4-0.start {
      counter-reset: lst-ctn-kix_list_4-0 0
    }

    .lst-kix_list_27-4 > li {
      counter-increment: lst-ctn-kix_list_27-4
    }

    .lst-kix_list_38-2 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) "." counter(lst-ctn-kix_list_38-2, decimal) " "
    }

    .lst-kix_list_38-4 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) "." counter(lst-ctn-kix_list_38-2, decimal) "." counter(lst-ctn-kix_list_38-3, decimal) "." counter(lst-ctn-kix_list_38-4, decimal) " "
    }

    ol.lst-kix_list_43-2.start {
      counter-reset: lst-ctn-kix_list_43-2 0
    }

    .lst-kix_list_16-1 > li {
      counter-increment: lst-ctn-kix_list_16-1
    }

    .lst-kix_list_27-1 > li {
      counter-increment: lst-ctn-kix_list_27-1
    }

    .lst-kix_list_38-4 > li {
      counter-increment: lst-ctn-kix_list_38-4
    }

    ol.lst-kix_list_43-4.start {
      counter-reset: lst-ctn-kix_list_43-4 0
    }

    .lst-kix_list_17-5 > li {
      counter-increment: lst-ctn-kix_list_17-5
    }

    ol.lst-kix_list_9-4.start {
      counter-reset: lst-ctn-kix_list_9-4 0
    }

    ol.lst-kix_list_41-8.start {
      counter-reset: lst-ctn-kix_list_41-8 0
    }

    .lst-kix_list_21-3 > li {
      counter-increment: lst-ctn-kix_list_21-3
    }

    ol.lst-kix_list_30-1.start {
      counter-reset: lst-ctn-kix_list_30-1 0
    }

    .lst-kix_list_36-4 > li:before {
      content: "" counter(lst-ctn-kix_list_36-4, lower-latin) ". "
    }

    ol.lst-kix_list_13-5.start {
      counter-reset: lst-ctn-kix_list_13-5 0
    }

    ol.lst-kix_list_13-8.start {
      counter-reset: lst-ctn-kix_list_13-8 0
    }

    .lst-kix_list_11-7 > li {
      counter-increment: lst-ctn-kix_list_11-7
    }

    .lst-kix_list_5-5 > li {
      counter-increment: lst-ctn-kix_list_5-5
    }

    .lst-kix_list_36-2 > li:before {
      content: "" counter(lst-ctn-kix_list_36-2, lower-roman) ". "
    }

    ol.lst-kix_list_43-7.start {
      counter-reset: lst-ctn-kix_list_43-7 0
    }

    .lst-kix_list_44-7 > li {
      counter-increment: lst-ctn-kix_list_44-7
    }

    .lst-kix_list_16-8 > li {
      counter-increment: lst-ctn-kix_list_16-8
    }

    .lst-kix_list_38-8 > li {
      counter-increment: lst-ctn-kix_list_38-8
    }

    ol.lst-kix_list_18-4.start {
      counter-reset: lst-ctn-kix_list_18-4 0
    }

    ol.lst-kix_list_29-1.start {
      counter-reset: lst-ctn-kix_list_29-1 0
    }

    ol.lst-kix_list_1-2.start {
      counter-reset: lst-ctn-kix_list_1-2 0
    }

    .lst-kix_list_20-4 > li:before {
      content: "" counter(lst-ctn-kix_list_20-4, decimal) ". "
    }

    ol.lst-kix_list_22-7.start {
      counter-reset: lst-ctn-kix_list_22-7 0
    }

    ol.lst-kix_list_6-1.start {
      counter-reset: lst-ctn-kix_list_6-1 0
    }

    .lst-kix_list_20-2 > li:before {
      content: "" counter(lst-ctn-kix_list_20-2, decimal) ". "
    }

    .lst-kix_list_54-1 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) " "
    }

    ol.lst-kix_list_16-8.start {
      counter-reset: lst-ctn-kix_list_16-8 0
    }

    ol.lst-kix_list_33-1.start {
      counter-reset: lst-ctn-kix_list_33-1 0
    }

    .lst-kix_list_4-8 > li {
      counter-increment: lst-ctn-kix_list_4-8
    }

    .lst-kix_list_55-5 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_9-1.start {
      counter-reset: lst-ctn-kix_list_9-1 0
    }

    ol.lst-kix_list_24-2.start {
      counter-reset: lst-ctn-kix_list_24-2 0
    }

    .lst-kix_list_54-7 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) "." counter(lst-ctn-kix_list_54-2, decimal) "." counter(lst-ctn-kix_list_54-3, decimal) "." counter(lst-ctn-kix_list_54-4, decimal) "." counter(lst-ctn-kix_list_54-5, decimal) "." counter(lst-ctn-kix_list_54-6, decimal) "." counter(lst-ctn-kix_list_54-7, decimal) " "
    }

    .lst-kix_list_19-8 > li {
      counter-increment: lst-ctn-kix_list_19-8
    }

    .lst-kix_list_55-3 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_35-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_35-6, decimal) ") "
    }

    .lst-kix_list_36-5 > li {
      counter-increment: lst-ctn-kix_list_36-5
    }

    .lst-kix_list_3-2 > li:before {
      content: "" counter(lst-ctn-kix_list_3-2, decimal) ". "
    }

    .lst-kix_list_8-1 > li:before {
      content: "o  "
    }

    .lst-kix_list_30-7 > li:before {
      content: "" counter(lst-ctn-kix_list_30-7, decimal) " "
    }

    .lst-kix_list_43-3 > li {
      counter-increment: lst-ctn-kix_list_43-3
    }

    .lst-kix_list_26-2 > li:before {
      content: "(" counter(lst-ctn-kix_list_26-2, lower-latin) ") "
    }

    .lst-kix_list_21-6 > li:before {
      content: "" counter(lst-ctn-kix_list_21-6, decimal) ". "
    }

    .lst-kix_list_41-3 > li {
      counter-increment: lst-ctn-kix_list_41-3
    }

    ol.lst-kix_list_54-6.start {
      counter-reset: lst-ctn-kix_list_54-6 0
    }

    ol.lst-kix_list_35-3.start {
      counter-reset: lst-ctn-kix_list_35-3 0
    }

    ol.lst-kix_list_4-2.start {
      counter-reset: lst-ctn-kix_list_4-2 0
    }

    ol.lst-kix_list_27-2.start {
      counter-reset: lst-ctn-kix_list_27-2 0
    }

    ol.lst-kix_list_18-7.start {
      counter-reset: lst-ctn-kix_list_18-7 0
    }

    .lst-kix_list_25-5 > li {
      counter-increment: lst-ctn-kix_list_25-5
    }

    ol.lst-kix_list_11-6.start {
      counter-reset: lst-ctn-kix_list_11-6 0
    }

    ol.lst-kix_list_6-4.start {
      counter-reset: lst-ctn-kix_list_6-4 0
    }

    .lst-kix_list_45-3 > li:before {
      content: "" counter(lst-ctn-kix_list_45-3, decimal) " "
    }

    .lst-kix_list_17-1 > li:before {
      content: "" counter(lst-ctn-kix_list_17-1, decimal) ". "
    }

    .lst-kix_list_32-3 > li {
      counter-increment: lst-ctn-kix_list_32-3
    }

    ol.lst-kix_list_27-1.start {
      counter-reset: lst-ctn-kix_list_27-1 0
    }

    .lst-kix_list_16-5 > li:before {
      content: "" counter(lst-ctn-kix_list_16-5, decimal) ". "
    }

    ol.lst-kix_list_22-2.start {
      counter-reset: lst-ctn-kix_list_22-2 0
    }

    .lst-kix_list_30-7 > li {
      counter-increment: lst-ctn-kix_list_30-7
    }

    ol.lst-kix_list_41-5.start {
      counter-reset: lst-ctn-kix_list_41-5 0
    }

    .lst-kix_list_44-7 > li:before {
      content: "" counter(lst-ctn-kix_list_44-7, decimal) " "
    }

    ol.lst-kix_list_29-3.start {
      counter-reset: lst-ctn-kix_list_29-3 0
    }

    .lst-kix_list_16-4 > li {
      counter-increment: lst-ctn-kix_list_16-4
    }

    ol.lst-kix_list_35-8.start {
      counter-reset: lst-ctn-kix_list_35-8 0
    }

    .lst-kix_list_38-1 > li {
      counter-increment: lst-ctn-kix_list_38-1
    }

    ol.lst-kix_list_16-5.start {
      counter-reset: lst-ctn-kix_list_16-5 0
    }

    .lst-kix_list_41-7 > li {
      counter-increment: lst-ctn-kix_list_41-7
    }

    ol.lst-kix_list_35-5.start {
      counter-reset: lst-ctn-kix_list_35-5 0
    }

    .lst-kix_list_2-6 > li:before {
      content: "" counter(lst-ctn-kix_list_2-6, decimal) ". "
    }

    ol.lst-kix_list_54-8.start {
      counter-reset: lst-ctn-kix_list_54-8 0
    }

    .lst-kix_list_14-5 > li {
      counter-increment: lst-ctn-kix_list_14-5
    }

    .lst-kix_list_7-5 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_27-6 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) "." counter(lst-ctn-kix_list_27-2, decimal) "." counter(lst-ctn-kix_list_27-3, decimal) "." counter(lst-ctn-kix_list_27-4, decimal) "." counter(lst-ctn-kix_list_27-5, decimal) "." counter(lst-ctn-kix_list_27-6, decimal) ". "
    }

    .lst-kix_list_23-2 > li {
      counter-increment: lst-ctn-kix_list_23-2
    }

    ol.lst-kix_list_30-4.start {
      counter-reset: lst-ctn-kix_list_30-4 0
    }

    .lst-kix_list_30-0 > li {
      counter-increment: lst-ctn-kix_list_30-0
    }

    ol.lst-kix_list_11-1.start {
      counter-reset: lst-ctn-kix_list_11-1 0
    }

    .lst-kix_list_18-5 > li:before {
      content: "" counter(lst-ctn-kix_list_18-5, decimal) ". "
    }

    .lst-kix_list_13-6 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) "." counter(lst-ctn-kix_list_13-2, decimal) "." counter(lst-ctn-kix_list_13-3, decimal) "." counter(lst-ctn-kix_list_13-4, decimal) "." counter(lst-ctn-kix_list_13-5, decimal) "." counter(lst-ctn-kix_list_13-6, decimal) ". "
    }

    .lst-kix_list_10-6 > li {
      counter-increment: lst-ctn-kix_list_10-6
    }

    .lst-kix_list_1-7 > li {
      counter-increment: lst-ctn-kix_list_1-7
    }

    ol.lst-kix_list_41-0.start {
      counter-reset: lst-ctn-kix_list_41-0 0
    }

    ol.lst-kix_list_54-5.start {
      counter-reset: lst-ctn-kix_list_54-5 0
    }

    .lst-kix_list_39-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_39-6, decimal) ") "
    }

    .lst-kix_list_54-6 > li {
      counter-increment: lst-ctn-kix_list_54-6
    }

    .lst-kix_list_29-0 > li {
      counter-increment: lst-ctn-kix_list_29-0
    }

    ol.lst-kix_list_24-7.start {
      counter-reset: lst-ctn-kix_list_24-7 0
    }

    .lst-kix_list_31-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_31-3, lower-latin) ") "
    }

    .lst-kix_list_43-6 > li {
      counter-increment: lst-ctn-kix_list_43-6
    }

    .lst-kix_list_10-2 > li:before {
      content: "(" counter(lst-ctn-kix_list_10-2, lower-roman) ") "
    }

    ol.lst-kix_list_1-7.start {
      counter-reset: lst-ctn-kix_list_1-7 0
    }

    .lst-kix_list_4-6 > li:before {
      content: "" counter(lst-ctn-kix_list_4-6, decimal) ". "
    }

    ol.lst-kix_list_29-8.start {
      counter-reset: lst-ctn-kix_list_29-8 0
    }

    .lst-kix_list_25-6 > li:before {
      content: "" counter(lst-ctn-kix_list_25-6, decimal) ". "
    }

    .lst-kix_list_46-7 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_34-2 > li {
      counter-increment: lst-ctn-kix_list_34-2
    }

    .lst-kix_list_12-2 > li {
      counter-increment: lst-ctn-kix_list_12-2
    }

    .lst-kix_list_9-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_9-5, upper-latin) ") "
    }

    .lst-kix_list_29-6 > li:before {
      content: "" counter(lst-ctn-kix_list_29-6, lower-roman) ") "
    }

    ol.lst-kix_list_22-0.start {
      counter-reset: lst-ctn-kix_list_22-0 0
    }

    ol.lst-kix_list_41-3.start {
      counter-reset: lst-ctn-kix_list_41-3 0
    }

    .lst-kix_list_12-2 > li:before {
      content: "(" counter(lst-ctn-kix_list_12-2, lower-roman) ") "
    }

    .lst-kix_list_11-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_11-6, decimal) ") "
    }

    .lst-kix_list_32-7 > li:before {
      content: "" counter(lst-ctn-kix_list_32-7, decimal) " "
    }

    ol.lst-kix_list_30-6.start {
      counter-reset: lst-ctn-kix_list_30-6 0
    }

    .lst-kix_list_1-2 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) "." counter(lst-ctn-kix_list_1-2, decimal) " "
    }

    ol.lst-kix_list_11-3.start {
      counter-reset: lst-ctn-kix_list_11-3 0
    }

    .lst-kix_list_49-3 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_1-0 > li {
      counter-increment: lst-ctn-kix_list_1-0
    }

    ol.lst-kix_list_41-2.start {
      counter-reset: lst-ctn-kix_list_41-2 0
    }

    .lst-kix_list_18-7 > li {
      counter-increment: lst-ctn-kix_list_18-7
    }

    .lst-kix_list_48-7 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_50-4 > li {
      counter-increment: lst-ctn-kix_list_50-4
    }

    .lst-kix_list_29-7 > li {
      counter-increment: lst-ctn-kix_list_29-7
    }

    ol.lst-kix_list_11-4.start {
      counter-reset: lst-ctn-kix_list_11-4 0
    }

    .lst-kix_list_28-2 > li:before {
      content: "" counter(lst-ctn-kix_list_28-2, lower-roman) ") "
    }

    .lst-kix_list_51-4 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) "." counter(lst-ctn-kix_list_51-2, decimal) "." counter(lst-ctn-kix_list_51-3, decimal) "." counter(lst-ctn-kix_list_51-4, decimal) " "
    }

    .lst-kix_list_14-1 > li:before {
      content: "" counter(lst-ctn-kix_list_14-1, lower-latin) ") "
    }

    .lst-kix_list_25-8 > li {
      counter-increment: lst-ctn-kix_list_25-8
    }

    .lst-kix_list_14-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_14-4, lower-latin) ") "
    }

    .lst-kix_list_51-1 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) " "
    }

    ol.lst-kix_list_45-4.start {
      counter-reset: lst-ctn-kix_list_45-4 0
    }

    .lst-kix_list_37-5 > li {
      counter-increment: lst-ctn-kix_list_37-5
    }

    .lst-kix_list_14-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_14-5, lower-roman) ") "
    }

    .lst-kix_list_14-7 > li:before {
      content: "" counter(lst-ctn-kix_list_14-7, lower-latin) ". "
    }

    .lst-kix_ek190fjkpfvz-7 > li:before {
      content: "\0025cb  "
    }

    .lst-kix_ek190fjkpfvz-6 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_34-6 > li {
      counter-increment: lst-ctn-kix_list_34-6
    }

    .lst-kix_list_51-8 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) "." counter(lst-ctn-kix_list_51-2, decimal) "." counter(lst-ctn-kix_list_51-3, decimal) "." counter(lst-ctn-kix_list_51-4, decimal) "." counter(lst-ctn-kix_list_51-5, decimal) "." counter(lst-ctn-kix_list_51-6, decimal) "." counter(lst-ctn-kix_list_51-7, decimal) "." counter(lst-ctn-kix_list_51-8, decimal) " "
    }

    .lst-kix_list_51-7 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) "." counter(lst-ctn-kix_list_51-2, decimal) "." counter(lst-ctn-kix_list_51-3, decimal) "." counter(lst-ctn-kix_list_51-4, decimal) "." counter(lst-ctn-kix_list_51-5, decimal) "." counter(lst-ctn-kix_list_51-6, decimal) "." counter(lst-ctn-kix_list_51-7, decimal) " "
    }

    .lst-kix_ek190fjkpfvz-0 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_32-1.start {
      counter-reset: lst-ctn-kix_list_32-1 0
    }

    .lst-kix_list_20-7 > li {
      counter-increment: lst-ctn-kix_list_20-7
    }

    ol.lst-kix_list_35-6.start {
      counter-reset: lst-ctn-kix_list_35-6 0
    }

    .lst-kix_ek190fjkpfvz-4 > li:before {
      content: "\0025cb  "
    }

    ol.lst-kix_list_28-3.start {
      counter-reset: lst-ctn-kix_list_28-3 0
    }

    .lst-kix_ek190fjkpfvz-3 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_44-1 > li {
      counter-increment: lst-ctn-kix_list_44-1
    }

    .lst-kix_list_32-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_32-3, lower-latin) ") "
    }

    .lst-kix_list_50-0 > li {
      counter-increment: lst-ctn-kix_list_50-0
    }

    .lst-kix_list_14-8 > li:before {
      content: "" counter(lst-ctn-kix_list_14-8, lower-roman) ". "
    }

    .lst-kix_list_39-7 > li {
      counter-increment: lst-ctn-kix_list_39-7
    }

    .lst-kix_list_32-0 > li:before {
      content: "" counter(lst-ctn-kix_list_32-0, decimal) " "
    }

    .lst-kix_list_3-2 > li {
      counter-increment: lst-ctn-kix_list_3-2
    }

    ol.lst-kix_list_15-5.start {
      counter-reset: lst-ctn-kix_list_15-5 0
    }

    ol.lst-kix_list_22-3.start {
      counter-reset: lst-ctn-kix_list_22-3 0
    }

    .lst-kix_list_5-4 > li {
      counter-increment: lst-ctn-kix_list_5-4
    }

    ol.lst-kix_list_24-6.start {
      counter-reset: lst-ctn-kix_list_24-6 0
    }

    .lst-kix_list_5-1 > li:before {
      content: "" counter(lst-ctn-kix_list_5-1, decimal) ". "
    }

    .lst-kix_list_5-7 > li:before {
      content: "" counter(lst-ctn-kix_list_5-7, decimal) ". "
    }

    .lst-kix_list_5-8 > li:before {
      content: "" counter(lst-ctn-kix_list_5-8, decimal) ". "
    }

    .lst-kix_list_5-4 > li:before {
      content: "" counter(lst-ctn-kix_list_5-4, decimal) ". "
    }

    .lst-kix_list_5-5 > li:before {
      content: "" counter(lst-ctn-kix_list_5-5, decimal) ". "
    }

    .lst-kix_list_50-1 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) " "
    }

    .lst-kix_list_50-2 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) "." counter(lst-ctn-kix_list_50-2, decimal) " "
    }

    .lst-kix_list_6-1 > li:before {
      content: "" counter(lst-ctn-kix_list_6-1, decimal) ". "
    }

    .lst-kix_list_6-3 > li:before {
      content: "" counter(lst-ctn-kix_list_6-3, decimal) ". "
    }

    ol.lst-kix_list_32-6.start {
      counter-reset: lst-ctn-kix_list_32-6 0
    }

    .lst-kix_list_6-8 > li {
      counter-increment: lst-ctn-kix_list_6-8
    }

    .lst-kix_list_6-0 > li:before {
      content: "" counter(lst-ctn-kix_list_6-0, decimal) ". "
    }

    .lst-kix_list_6-4 > li:before {
      content: "" counter(lst-ctn-kix_list_6-4, decimal) ". "
    }

    ol.lst-kix_list_14-8.start {
      counter-reset: lst-ctn-kix_list_14-8 0
    }

    ol.lst-kix_list_15-0.start {
      counter-reset: lst-ctn-kix_list_15-0 0
    }

    ol.lst-kix_list_44-7.start {
      counter-reset: lst-ctn-kix_list_44-7 0
    }

    .lst-kix_list_2-5 > li {
      counter-increment: lst-ctn-kix_list_2-5
    }

    .lst-kix_list_50-8 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) "." counter(lst-ctn-kix_list_50-2, decimal) "." counter(lst-ctn-kix_list_50-3, decimal) "." counter(lst-ctn-kix_list_50-4, decimal) "." counter(lst-ctn-kix_list_50-5, decimal) "." counter(lst-ctn-kix_list_50-6, decimal) "." counter(lst-ctn-kix_list_50-7, decimal) "." counter(lst-ctn-kix_list_50-8, decimal) " "
    }

    .lst-kix_list_50-5 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) "." counter(lst-ctn-kix_list_50-2, decimal) "." counter(lst-ctn-kix_list_50-3, decimal) "." counter(lst-ctn-kix_list_50-4, decimal) "." counter(lst-ctn-kix_list_50-5, decimal) " "
    }

    .lst-kix_list_6-7 > li:before {
      content: "" counter(lst-ctn-kix_list_6-7, decimal) ". "
    }

    .lst-kix_list_6-6 > li:before {
      content: "" counter(lst-ctn-kix_list_6-6, decimal) ". "
    }

    ol.lst-kix_list_10-6.start {
      counter-reset: lst-ctn-kix_list_10-6 0
    }

    ol.lst-kix_list_27-6.start {
      counter-reset: lst-ctn-kix_list_27-6 0
    }

    .lst-kix_list_7-6 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_18-5 > li {
      counter-increment: lst-ctn-kix_list_18-5
    }

    ol.lst-kix_list_19-7.start {
      counter-reset: lst-ctn-kix_list_19-7 0
    }

    ol.lst-kix_list_6-2.start {
      counter-reset: lst-ctn-kix_list_6-2 0
    }

    .lst-kix_list_15-5 > li {
      counter-increment: lst-ctn-kix_list_15-5
    }

    ol.lst-kix_list_24-1.start {
      counter-reset: lst-ctn-kix_list_24-1 0
    }

    ol.lst-kix_list_36-3.start {
      counter-reset: lst-ctn-kix_list_36-3 0
    }

    .lst-kix_list_7-2 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_27-2 > li {
      counter-increment: lst-ctn-kix_list_27-2
    }

    .lst-kix_list_31-0 > li {
      counter-increment: lst-ctn-kix_list_31-0
    }

    ol.lst-kix_list_22-8.start {
      counter-reset: lst-ctn-kix_list_22-8 0
    }

    .lst-kix_list_34-8 > li:before {
      content: "" counter(lst-ctn-kix_list_34-8, decimal) " "
    }

    .lst-kix_list_31-0 > li:before {
      content: "" counter(lst-ctn-kix_list_31-0, decimal) " "
    }

    .lst-kix_list_12-6 > li {
      counter-increment: lst-ctn-kix_list_12-6
    }

    ol.lst-kix_list_23-0.start {
      counter-reset: lst-ctn-kix_list_23-0 0
    }

    .lst-kix_list_52-7 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_9-8 > li {
      counter-increment: lst-ctn-kix_list_9-8
    }

    .lst-kix_list_13-4 > li {
      counter-increment: lst-ctn-kix_list_13-4
    }

    .lst-kix_list_52-3 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_31-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_31-4, lower-roman) ") "
    }

    .lst-kix_list_15-7 > li:before {
      content: "" counter(lst-ctn-kix_list_15-7, lower-latin) ". "
    }

    .lst-kix_list_17-7 > li {
      counter-increment: lst-ctn-kix_list_17-7
    }

    .lst-kix_list_4-5 > li:before {
      content: "" counter(lst-ctn-kix_list_4-5, decimal) ". "
    }

    .lst-kix_list_51-4 > li {
      counter-increment: lst-ctn-kix_list_51-4
    }

    .lst-kix_list_15-1 > li:before {
      content: "Section " counter(lst-ctn-kix_list_15-0, upper-roman) "." counter(lst-ctn-kix_list_15-1, decimal) " "
    }

    ol.lst-kix_list_1-4.start {
      counter-reset: lst-ctn-kix_list_1-4 0
    }

    .lst-kix_list_15-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_15-3, lower-roman) ") "
    }

    .lst-kix_list_22-1 > li {
      counter-increment: lst-ctn-kix_list_22-1
    }

    .lst-kix_list_42-7 > li {
      counter-increment: lst-ctn-kix_list_42-7
    }

    ol.lst-kix_list_4-4.start {
      counter-reset: lst-ctn-kix_list_4-4 0
    }

    .lst-kix_list_33-2 > li {
      counter-increment: lst-ctn-kix_list_33-2
    }

    .lst-kix_list_40-5 > li {
      counter-increment: lst-ctn-kix_list_40-5
    }

    .lst-kix_list_30-2 > li {
      counter-increment: lst-ctn-kix_list_30-2
    }

    ol.lst-kix_list_9-2.start {
      counter-reset: lst-ctn-kix_list_9-2 0
    }

    .lst-kix_list_53-5 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_11-2 > li {
      counter-increment: lst-ctn-kix_list_11-2
    }

    .lst-kix_list_33-8 > li {
      counter-increment: lst-ctn-kix_list_33-8
    }

    .lst-kix_list_12-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_12-3, upper-latin) ") "
    }

    ol.lst-kix_list_31-4.start {
      counter-reset: lst-ctn-kix_list_31-4 0
    }

    .lst-kix_list_32-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_32-6, decimal) ") "
    }

    ol.lst-kix_list_54-0.start {
      counter-reset: lst-ctn-kix_list_54-0 6
    }

    .lst-kix_list_33-2 > li:before {
      content: "" counter(lst-ctn-kix_list_33-0, decimal) "." counter(lst-ctn-kix_list_33-1, decimal) "." counter(lst-ctn-kix_list_33-2, decimal) " "
    }

    .lst-kix_list_16-3 > li {
      counter-increment: lst-ctn-kix_list_16-3
    }

    .lst-kix_list_13-3 > li {
      counter-increment: lst-ctn-kix_list_13-3
    }

    ol.lst-kix_list_40-5.start {
      counter-reset: lst-ctn-kix_list_40-5 0
    }

    .lst-kix_list_10-4 > li {
      counter-increment: lst-ctn-kix_list_10-4
    }

    .lst-kix_list_14-1 > li {
      counter-increment: lst-ctn-kix_list_14-1
    }

    .lst-kix_list_34-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_34-4, lower-roman) ") "
    }

    .lst-kix_list_35-3 > li {
      counter-increment: lst-ctn-kix_list_35-3
    }

    .lst-kix_list_28-0 > li {
      counter-increment: lst-ctn-kix_list_28-0
    }

    .lst-kix_list_13-5 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) "." counter(lst-ctn-kix_list_13-2, decimal) "." counter(lst-ctn-kix_list_13-3, decimal) "." counter(lst-ctn-kix_list_13-4, decimal) "." counter(lst-ctn-kix_list_13-5, decimal) ". "
    }

    ol.lst-kix_list_28-8.start {
      counter-reset: lst-ctn-kix_list_28-8 0
    }

    .lst-kix_list_36-1 > li {
      counter-increment: lst-ctn-kix_list_36-1
    }

    .lst-kix_list_36-7 > li {
      counter-increment: lst-ctn-kix_list_36-7
    }

    .lst-kix_list_29-4 > li {
      counter-increment: lst-ctn-kix_list_29-4
    }

    .lst-kix_list_54-4 > li {
      counter-increment: lst-ctn-kix_list_54-4
    }

    .lst-kix_list_33-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_33-6, decimal) ") "
    }

    .lst-kix_list_32-4 > li {
      counter-increment: lst-ctn-kix_list_32-4
    }

    .lst-kix_list_33-8 > li:before {
      content: "" counter(lst-ctn-kix_list_33-8, decimal) " "
    }

    ol.lst-kix_list_29-0.start {
      counter-reset: lst-ctn-kix_list_29-0 0
    }

    .lst-kix_list_34-2 > li:before {
      content: "" counter(lst-ctn-kix_list_34-0, decimal) "." counter(lst-ctn-kix_list_34-1, decimal) "." counter(lst-ctn-kix_list_34-2, decimal) " "
    }

    ol.lst-kix_list_54-1.start {
      counter-reset: lst-ctn-kix_list_54-1 0
    }

    .lst-kix_list_34-5 > li {
      counter-increment: lst-ctn-kix_list_34-5
    }

    .lst-kix_list_30-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_30-5, upper-latin) ") "
    }

    .lst-kix_list_35-0 > li:before {
      content: "" counter(lst-ctn-kix_list_35-0, decimal) " "
    }

    .lst-kix_list_35-1 > li:before {
      content: "" counter(lst-ctn-kix_list_35-1, decimal) ". "
    }

    .lst-kix_list_35-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_35-4, lower-roman) ") "
    }

    ol.lst-kix_list_40-6.start {
      counter-reset: lst-ctn-kix_list_40-6 0
    }

    .lst-kix_list_28-8 > li {
      counter-increment: lst-ctn-kix_list_28-8
    }

    .lst-kix_list_1-1 > li {
      counter-increment: lst-ctn-kix_list_1-1
    }

    .lst-kix_list_30-1 > li:before {
      content: "" counter(lst-ctn-kix_list_30-1, decimal) ". "
    }

    .lst-kix_list_3-0 > li:before {
      content: "" counter(lst-ctn-kix_list_3-0, decimal) ". "
    }

    .lst-kix_list_30-2 > li:before {
      content: "" counter(lst-ctn-kix_list_30-0, decimal) "." counter(lst-ctn-kix_list_30-1, decimal) "." counter(lst-ctn-kix_list_30-2, decimal) " "
    }

    .lst-kix_list_4-0 > li {
      counter-increment: lst-ctn-kix_list_4-0
    }

    .lst-kix_list_31-6 > li {
      counter-increment: lst-ctn-kix_list_31-6
    }

    ol.lst-kix_list_37-5.start {
      counter-reset: lst-ctn-kix_list_37-5 0
    }

    .lst-kix_list_3-3 > li:before {
      content: "" counter(lst-ctn-kix_list_3-3, decimal) ". "
    }

    .lst-kix_list_38-3 > li {
      counter-increment: lst-ctn-kix_list_38-3
    }

    ol.lst-kix_list_10-7.start {
      counter-reset: lst-ctn-kix_list_10-7 0
    }

    .lst-kix_list_53-2 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_53-3 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_17-1 > li {
      counter-increment: lst-ctn-kix_list_17-1
    }

    .lst-kix_list_11-1 > li:before {
      content: "" counter(lst-ctn-kix_list_11-0, decimal) "." counter(lst-ctn-kix_list_11-1, decimal) ". "
    }

    ol.lst-kix_list_31-0.start {
      counter-reset: lst-ctn-kix_list_31-0 0
    }

    .lst-kix_list_11-0 > li:before {
      content: "" counter(lst-ctn-kix_list_11-0, decimal) " "
    }

    ol.lst-kix_list_9-3.start {
      counter-reset: lst-ctn-kix_list_9-3 0
    }

    .lst-kix_list_8-8 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_35-4 > li {
      counter-increment: lst-ctn-kix_list_35-4
    }

    .lst-kix_list_43-5 > li {
      counter-increment: lst-ctn-kix_list_43-5
    }

    .lst-kix_list_4-8 > li:before {
      content: "" counter(lst-ctn-kix_list_4-8, decimal) ". "
    }

    .lst-kix_list_21-5 > li {
      counter-increment: lst-ctn-kix_list_21-5
    }

    .lst-kix_list_14-2 > li {
      counter-increment: lst-ctn-kix_list_14-2
    }

    .lst-kix_list_16-0 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_40-6 > li {
      counter-increment: lst-ctn-kix_list_40-6
    }

    ol.lst-kix_list_4-8.start {
      counter-reset: lst-ctn-kix_list_4-8 0
    }

    .lst-kix_list_16-3 > li:before {
      content: "" counter(lst-ctn-kix_list_16-3, decimal) ". "
    }

    ol.lst-kix_list_37-4.start {
      counter-reset: lst-ctn-kix_list_37-4 0
    }

    .lst-kix_list_41-1 > li {
      counter-increment: lst-ctn-kix_list_41-1
    }

    .lst-kix_list_44-0 > li {
      counter-increment: lst-ctn-kix_list_44-0
    }

    .lst-kix_list_17-8 > li:before {
      content: "" counter(lst-ctn-kix_list_17-8, decimal) ". "
    }

    .lst-kix_list_7-0 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_27-7.start {
      counter-reset: lst-ctn-kix_list_27-7 0
    }

    ol.lst-kix_list_19-6.start {
      counter-reset: lst-ctn-kix_list_19-6 0
    }

    ol.lst-kix_list_9-7.start {
      counter-reset: lst-ctn-kix_list_9-7 0
    }

    .lst-kix_list_2-4 > li:before {
      content: "" counter(lst-ctn-kix_list_2-4, decimal) ". "
    }

    .lst-kix_list_7-3 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_48-5 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_9-7 > li {
      counter-increment: lst-ctn-kix_list_9-7
    }

    .lst-kix_list_13-8 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) "." counter(lst-ctn-kix_list_13-2, decimal) "." counter(lst-ctn-kix_list_13-3, decimal) "." counter(lst-ctn-kix_list_13-4, decimal) "." counter(lst-ctn-kix_list_13-5, decimal) "." counter(lst-ctn-kix_list_13-6, decimal) "." counter(lst-ctn-kix_list_13-7, decimal) "." counter(lst-ctn-kix_list_13-8, decimal) ". "
    }

    .lst-kix_list_31-1 > li:before {
      content: "" counter(lst-ctn-kix_list_31-1, decimal) ". "
    }

    .lst-kix_list_18-7 > li:before {
      content: "" counter(lst-ctn-kix_list_18-7, decimal) ". "
    }

    .lst-kix_list_48-1 > li:before {
      content: "o  "
    }

    .lst-kix_list_52-6 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_15-4 > li:before {
      content: "" counter(lst-ctn-kix_list_15-4, decimal) ") "
    }

    ol.lst-kix_list_19-1.start {
      counter-reset: lst-ctn-kix_list_19-1 0
    }

    .lst-kix_list_10-4 > li:before {
      content: "" counter(lst-ctn-kix_list_10-4, decimal) " "
    }

    .lst-kix_list_10-8 > li:before {
      content: "" counter(lst-ctn-kix_list_10-8, decimal) " "
    }

    .lst-kix_list_15-0 > li:before {
      content: "Article " counter(lst-ctn-kix_list_15-0, upper-roman) ". "
    }

    ol.lst-kix_list_14-3.start {
      counter-reset: lst-ctn-kix_list_14-3 0
    }

    ol.lst-kix_list_44-2.start {
      counter-reset: lst-ctn-kix_list_44-2 0
    }

    ol.lst-kix_list_32-5.start {
      counter-reset: lst-ctn-kix_list_32-5 0
    }

    .lst-kix_list_12-8 > li {
      counter-increment: lst-ctn-kix_list_12-8
    }

    ol.lst-kix_list_14-4.start {
      counter-reset: lst-ctn-kix_list_14-4 0
    }

    .lst-kix_list_9-7 > li:before {
      content: "" counter(lst-ctn-kix_list_9-7, decimal) " "
    }

    .lst-kix_list_2-4 > li {
      counter-increment: lst-ctn-kix_list_2-4
    }

    .lst-kix_list_29-4 > li:before {
      content: "" counter(lst-ctn-kix_list_29-4, decimal) ") "
    }

    .lst-kix_list_53-6 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_32-2.start {
      counter-reset: lst-ctn-kix_list_32-2 0
    }

    .lst-kix_list_11-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_11-4, lower-roman) ") "
    }

    ol.lst-kix_list_37-0.start {
      counter-reset: lst-ctn-kix_list_37-0 0
    }

    .lst-kix_list_29-0 > li:before {
      content: "Article " counter(lst-ctn-kix_list_29-0, upper-roman) ". "
    }

    ol.lst-kix_list_19-2.start {
      counter-reset: lst-ctn-kix_list_19-2 0
    }

    .lst-kix_list_54-2 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) "." counter(lst-ctn-kix_list_54-2, decimal) " "
    }

    .lst-kix_list_12-0 > li:before {
      content: "" counter(lst-ctn-kix_list_12-0, decimal) " "
    }

    .lst-kix_list_1-4 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) "." counter(lst-ctn-kix_list_1-2, decimal) "." counter(lst-ctn-kix_list_1-3, decimal) "." counter(lst-ctn-kix_list_1-4, decimal) " "
    }

    ol.lst-kix_list_44-3.start {
      counter-reset: lst-ctn-kix_list_44-3 0
    }

    .lst-kix_list_1-6 > li {
      counter-increment: lst-ctn-kix_list_1-6
    }

    .lst-kix_list_34-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_34-5, upper-latin) ") "
    }

    ol.lst-kix_list_45-8.start {
      counter-reset: lst-ctn-kix_list_45-8 0
    }

    .lst-kix_list_33-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_33-5, upper-latin) ") "
    }

    ol.lst-kix_list_19-3.start {
      counter-reset: lst-ctn-kix_list_19-3 0
    }

    .lst-kix_list_2-0 > li:before {
      content: "" counter(lst-ctn-kix_list_2-0, decimal) ". "
    }

    ol.lst-kix_list_9-8.start {
      counter-reset: lst-ctn-kix_list_9-8 0
    }

    .lst-kix_list_1-8 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) "." counter(lst-ctn-kix_list_1-2, decimal) "." counter(lst-ctn-kix_list_1-3, decimal) "." counter(lst-ctn-kix_list_1-4, decimal) "." counter(lst-ctn-kix_list_1-5, decimal) "." counter(lst-ctn-kix_list_1-6, decimal) "." counter(lst-ctn-kix_list_1-7, decimal) "." counter(lst-ctn-kix_list_1-8, decimal) " "
    }

    .lst-kix_list_49-1 > li:before {
      content: "o  "
    }

    .lst-kix_list_34-1 > li:before {
      content: "" counter(lst-ctn-kix_list_34-1, decimal) ". "
    }

    .lst-kix_list_19-0 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_44-6.start {
      counter-reset: lst-ctn-kix_list_44-6 0
    }

    .lst-kix_list_29-6 > li {
      counter-increment: lst-ctn-kix_list_29-6
    }

    .lst-kix_list_47-2 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_31-3.start {
      counter-reset: lst-ctn-kix_list_31-3 0
    }

    ol.lst-kix_list_19-8.start {
      counter-reset: lst-ctn-kix_list_19-8 0
    }

    ol.lst-kix_list_36-4.start {
      counter-reset: lst-ctn-kix_list_36-4 0
    }

    ol.lst-kix_list_14-7.start {
      counter-reset: lst-ctn-kix_list_14-7 0
    }

    .lst-kix_list_19-7 > li:before {
      content: "" counter(lst-ctn-kix_list_19-7, decimal) ". "
    }

    .lst-kix_list_9-2 > li {
      counter-increment: lst-ctn-kix_list_9-2
    }

    ol.lst-kix_list_23-8.start {
      counter-reset: lst-ctn-kix_list_23-8 0
    }

    .lst-kix_list_24-5 > li {
      counter-increment: lst-ctn-kix_list_24-5
    }

    .lst-kix_list_46-3 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_38-1 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) " "
    }

    ol.lst-kix_list_15-6.start {
      counter-reset: lst-ctn-kix_list_15-6 0
    }

    .lst-kix_list_37-2 > li:before {
      content: "" counter(lst-ctn-kix_list_37-0, decimal) "." counter(lst-ctn-kix_list_37-1, decimal) "." counter(lst-ctn-kix_list_37-2, decimal) " "
    }

    .lst-kix_list_25-6 > li {
      counter-increment: lst-ctn-kix_list_25-6
    }

    .lst-kix_list_37-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_37-5, upper-latin) ") "
    }

    .lst-kix_list_46-0 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_31-8.start {
      counter-reset: lst-ctn-kix_list_31-8 0
    }

    .lst-kix_list_18-1 > li:before {
      content: "" counter(lst-ctn-kix_list_18-1, decimal) ". "
    }

    .lst-kix_list_38-8 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) "." counter(lst-ctn-kix_list_38-2, decimal) "." counter(lst-ctn-kix_list_38-3, decimal) "." counter(lst-ctn-kix_list_38-4, decimal) "." counter(lst-ctn-kix_list_38-5, decimal) "." counter(lst-ctn-kix_list_38-6, decimal) "." counter(lst-ctn-kix_list_38-7, decimal) "." counter(lst-ctn-kix_list_38-8, decimal) " "
    }

    .lst-kix_list_40-3 > li {
      counter-increment: lst-ctn-kix_list_40-3
    }

    .lst-kix_list_44-3 > li {
      counter-increment: lst-ctn-kix_list_44-3
    }

    .lst-kix_list_45-4 > li {
      counter-increment: lst-ctn-kix_list_45-4
    }

    .lst-kix_list_23-4 > li {
      counter-increment: lst-ctn-kix_list_23-4
    }

    ol.lst-kix_list_23-1.start {
      counter-reset: lst-ctn-kix_list_23-1 0
    }

    ol.lst-kix_list_45-5.start {
      counter-reset: lst-ctn-kix_list_45-5 0
    }

    ol.lst-kix_list_32-0.start {
      counter-reset: lst-ctn-kix_list_32-0 0
    }

    .lst-kix_list_23-1 > li {
      counter-increment: lst-ctn-kix_list_23-1
    }

    .lst-kix_list_2-7 > li:before {
      content: "" counter(lst-ctn-kix_list_2-7, decimal) ". "
    }

    .lst-kix_list_2-7 > li {
      counter-increment: lst-ctn-kix_list_2-7
    }

    .lst-kix_list_24-2 > li {
      counter-increment: lst-ctn-kix_list_24-2
    }

    .lst-kix_list_27-5 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) "." counter(lst-ctn-kix_list_27-2, decimal) "." counter(lst-ctn-kix_list_27-3, decimal) "." counter(lst-ctn-kix_list_27-4, decimal) "." counter(lst-ctn-kix_list_27-5, decimal) ". "
    }

    ol.lst-kix_list_10-3.start {
      counter-reset: lst-ctn-kix_list_10-3 0
    }

    .lst-kix_list_22-3 > li {
      counter-increment: lst-ctn-kix_list_22-3
    }

    .lst-kix_list_39-7 > li:before {
      content: "" counter(lst-ctn-kix_list_39-7, decimal) " "
    }

    .lst-kix_list_55-7 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_10-1 > li:before {
      content: "(" counter(lst-ctn-kix_list_10-1, lower-latin) ") "
    }

    .lst-kix_list_18-4 > li:before {
      content: "" counter(lst-ctn-kix_list_18-4, decimal) ". "
    }

    ol.lst-kix_list_15-1.start {
      counter-reset: lst-ctn-kix_list_15-1 0
    }

    ol.lst-kix_list_15-4.start {
      counter-reset: lst-ctn-kix_list_15-4 0
    }

    ol.lst-kix_list_45-3.start {
      counter-reset: lst-ctn-kix_list_45-3 0
    }

    .lst-kix_list_36-6 > li:before {
      content: "" counter(lst-ctn-kix_list_36-6, decimal) ". "
    }

    .lst-kix_list_50-6 > li {
      counter-increment: lst-ctn-kix_list_50-6
    }

    .lst-kix_list_36-0 > li:before {
      content: "" counter(lst-ctn-kix_list_36-0, upper-latin) " "
    }

    ol.lst-kix_list_40-2.start {
      counter-reset: lst-ctn-kix_list_40-2 0
    }

    ol.lst-kix_list_44-8.start {
      counter-reset: lst-ctn-kix_list_44-8 0
    }

    .lst-kix_list_34-0 > li {
      counter-increment: lst-ctn-kix_list_34-0
    }

    ol.lst-kix_list_45-0.start {
      counter-reset: lst-ctn-kix_list_45-0 0
    }

    .lst-kix_list_26-4 > li {
      counter-increment: lst-ctn-kix_list_26-4
    }

    ol.lst-kix_list_5-7.start {
      counter-reset: lst-ctn-kix_list_5-7 0
    }

    ol.lst-kix_list_28-7.start {
      counter-reset: lst-ctn-kix_list_28-7 0
    }

    .lst-kix_list_20-8 > li:before {
      content: "" counter(lst-ctn-kix_list_20-8, decimal) ". "
    }

    .lst-kix_list_39-1 > li {
      counter-increment: lst-ctn-kix_list_39-1
    }

    .lst-kix_list_46-6 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_29-7 > li:before {
      content: "" counter(lst-ctn-kix_list_29-7, lower-latin) ". "
    }

    .lst-kix_list_9-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_9-4, lower-roman) ") "
    }

    .lst-kix_list_30-8 > li {
      counter-increment: lst-ctn-kix_list_30-8
    }

    ol.lst-kix_list_36-2.start {
      counter-reset: lst-ctn-kix_list_36-2 0
    }

    .lst-kix_list_1-1 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) " "
    }

    .lst-kix_list_11-7 > li:before {
      content: "" counter(lst-ctn-kix_list_11-7, decimal) " "
    }

    .lst-kix_list_49-4 > li:before {
      content: "o  "
    }

    .lst-kix_list_55-1 > li:before {
      content: "o  "
    }

    ol.lst-kix_list_10-5.start {
      counter-reset: lst-ctn-kix_list_10-5 0
    }

    .lst-kix_list_54-5 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) "." counter(lst-ctn-kix_list_54-2, decimal) "." counter(lst-ctn-kix_list_54-3, decimal) "." counter(lst-ctn-kix_list_54-4, decimal) "." counter(lst-ctn-kix_list_54-5, decimal) " "
    }

    .lst-kix_list_14-7 > li {
      counter-increment: lst-ctn-kix_list_14-7
    }

    ol.lst-kix_list_31-5.start {
      counter-reset: lst-ctn-kix_list_31-5 0
    }

    .lst-kix_list_48-8 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_40-4.start {
      counter-reset: lst-ctn-kix_list_40-4 0
    }

    .lst-kix_list_28-1 > li:before {
      content: "" counter(lst-ctn-kix_list_28-1, lower-latin) ") "
    }

    .lst-kix_list_30-8 > li:before {
      content: "" counter(lst-ctn-kix_list_30-8, decimal) " "
    }

    .lst-kix_list_35-7 > li:before {
      content: "" counter(lst-ctn-kix_list_35-7, decimal) " "
    }

    .lst-kix_list_30-5 > li {
      counter-increment: lst-ctn-kix_list_30-5
    }

    .lst-kix_list_26-6 > li:before {
      content: "" counter(lst-ctn-kix_list_26-6, lower-roman) ") "
    }

    .lst-kix_list_8-2 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_36-7.start {
      counter-reset: lst-ctn-kix_list_36-7 0
    }

    ol.lst-kix_list_31-6.start {
      counter-reset: lst-ctn-kix_list_31-6 0
    }

    .lst-kix_list_12-0 > li {
      counter-increment: lst-ctn-kix_list_12-0
    }

    ol.lst-kix_list_45-1.start {
      counter-reset: lst-ctn-kix_list_45-1 0
    }

    ol.lst-kix_list_40-0.start {
      counter-reset: lst-ctn-kix_list_40-0 0
    }

    .lst-kix_list_8-5 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_28-0.start {
      counter-reset: lst-ctn-kix_list_28-0 0
    }

    .lst-kix_list_26-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_26-3, lower-roman) ") "
    }

    .lst-kix_list_3-6 > li:before {
      content: "" counter(lst-ctn-kix_list_3-6, decimal) ". "
    }

    .lst-kix_list_21-7 > li:before {
      content: "" counter(lst-ctn-kix_list_21-7, decimal) ". "
    }

    ol.lst-kix_list_5-0.start {
      counter-reset: lst-ctn-kix_list_5-0 0
    }

    .lst-kix_list_39-4 > li {
      counter-increment: lst-ctn-kix_list_39-4
    }

    ol.lst-kix_list_31-7.start {
      counter-reset: lst-ctn-kix_list_31-7 0
    }

    .lst-kix_list_16-6 > li:before {
      content: "" counter(lst-ctn-kix_list_16-6, decimal) ". "
    }

    ol.lst-kix_list_10-1.start {
      counter-reset: lst-ctn-kix_list_10-1 0
    }

    .lst-kix_list_18-2 > li {
      counter-increment: lst-ctn-kix_list_18-2
    }

    ol.lst-kix_list_5-6.start {
      counter-reset: lst-ctn-kix_list_5-6 0
    }

    .lst-kix_list_22-6 > li {
      counter-increment: lst-ctn-kix_list_22-6
    }

    .lst-kix_list_25-0 > li {
      counter-increment: lst-ctn-kix_list_25-0
    }

    .lst-kix_list_45-4 > li:before {
      content: "" counter(lst-ctn-kix_list_45-4, decimal) " "
    }

    .lst-kix_list_45-7 > li:before {
      content: "" counter(lst-ctn-kix_list_45-7, decimal) " "
    }

    ol.lst-kix_list_28-6.start {
      counter-reset: lst-ctn-kix_list_28-6 0
    }

    .lst-kix_list_19-3 > li {
      counter-increment: lst-ctn-kix_list_19-3
    }

    ol.lst-kix_list_28-5.start {
      counter-reset: lst-ctn-kix_list_28-5 0
    }

    .lst-kix_list_44-8 > li:before {
      content: "" counter(lst-ctn-kix_list_44-8, decimal) " "
    }

    .lst-kix_list_23-7 > li {
      counter-increment: lst-ctn-kix_list_23-7
    }

    .lst-kix_list_44-6 > li {
      counter-increment: lst-ctn-kix_list_44-6
    }

    ol.lst-kix_list_10-2.start {
      counter-reset: lst-ctn-kix_list_10-2 0
    }

    .lst-kix_list_40-0 > li {
      counter-increment: lst-ctn-kix_list_40-0
    }

    .lst-kix_list_17-2 > li:before {
      content: "" counter(lst-ctn-kix_list_17-2, decimal) ". "
    }

    ol.lst-kix_list_5-5.start {
      counter-reset: lst-ctn-kix_list_5-5 0
    }

    ol.lst-kix_list_40-1.start {
      counter-reset: lst-ctn-kix_list_40-1 0
    }

    .lst-kix_list_45-7 > li {
      counter-increment: lst-ctn-kix_list_45-7
    }

    .lst-kix_list_17-5 > li:before {
      content: "" counter(lst-ctn-kix_list_17-5, decimal) ". "
    }

    .lst-kix_list_6-2 > li {
      counter-increment: lst-ctn-kix_list_6-2
    }

    ol.lst-kix_list_36-6.start {
      counter-reset: lst-ctn-kix_list_36-6 0
    }

    .lst-kix_list_27-2 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) "." counter(lst-ctn-kix_list_27-2, decimal) ". "
    }

    .lst-kix_list_22-3 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) "." counter(lst-ctn-kix_list_22-2, decimal) "." counter(lst-ctn-kix_list_22-3, decimal) ". "
    }

    ol.lst-kix_list_28-4.start {
      counter-reset: lst-ctn-kix_list_28-4 0
    }

    .lst-kix_list_43-8 > li {
      counter-increment: lst-ctn-kix_list_43-8
    }

    ol.lst-kix_list_23-3.start {
      counter-reset: lst-ctn-kix_list_23-3 0
    }

    ol.lst-kix_list_5-4.start {
      counter-reset: lst-ctn-kix_list_5-4 0
    }

    .lst-kix_list_52-0 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_5-1.start {
      counter-reset: lst-ctn-kix_list_5-1 0
    }

    .lst-kix_list_25-3 > li {
      counter-increment: lst-ctn-kix_list_25-3
    }

    .lst-kix_list_16-6 > li {
      counter-increment: lst-ctn-kix_list_16-6
    }

    ol.lst-kix_list_28-1.start {
      counter-reset: lst-ctn-kix_list_28-1 0
    }

    .lst-kix_list_31-7 > li:before {
      content: "" counter(lst-ctn-kix_list_31-7, decimal) " "
    }

    ol.lst-kix_list_10-0.start {
      counter-reset: lst-ctn-kix_list_10-0 0
    }

    .lst-kix_list_3-8 > li {
      counter-increment: lst-ctn-kix_list_3-8
    }

    .lst-kix_list_4-6 > li {
      counter-increment: lst-ctn-kix_list_4-6
    }

    .lst-kix_list_51-7 > li {
      counter-increment: lst-ctn-kix_list_51-7
    }

    .lst-kix_list_4-2 > li:before {
      content: "" counter(lst-ctn-kix_list_4-2, decimal) ". "
    }

    ol.lst-kix_list_23-6.start {
      counter-reset: lst-ctn-kix_list_23-6 0
    }

    .lst-kix_list_17-4 > li {
      counter-increment: lst-ctn-kix_list_17-4
    }

    .lst-kix_list_36-3 > li:before {
      content: "" counter(lst-ctn-kix_list_36-3, decimal) ". "
    }

    .lst-kix_list_26-1 > li {
      counter-increment: lst-ctn-kix_list_26-1
    }

    .lst-kix_list_9-1 > li:before {
      content: "" counter(lst-ctn-kix_list_9-1, decimal) ". "
    }

    .lst-kix_list_15-8 > li {
      counter-increment: lst-ctn-kix_list_15-8
    }

    .lst-kix_list_40-8 > li:before {
      content: "" counter(lst-ctn-kix_list_40-8, decimal) " "
    }

    ol.lst-kix_list_36-8.start {
      counter-reset: lst-ctn-kix_list_36-8 0
    }

    .lst-kix_list_54-1 > li {
      counter-increment: lst-ctn-kix_list_54-1
    }

    .lst-kix_list_37-8 > li {
      counter-increment: lst-ctn-kix_list_37-8
    }

    .lst-kix_list_31-3 > li {
      counter-increment: lst-ctn-kix_list_31-3
    }

    .lst-kix_list_41-4 > li:before {
      content: "" counter(lst-ctn-kix_list_41-4, decimal) " "
    }

    ol.lst-kix_list_23-5.start {
      counter-reset: lst-ctn-kix_list_23-5 0
    }

    .lst-kix_list_10-1 > li {
      counter-increment: lst-ctn-kix_list_10-1
    }

    .lst-kix_list_49-7 > li:before {
      content: "o  "
    }

    ol.lst-kix_list_28-2.start {
      counter-reset: lst-ctn-kix_list_28-2 0
    }

    .lst-kix_list_12-6 > li:before {
      content: "" counter(lst-ctn-kix_list_12-6, decimal) ". "
    }

    ol.lst-kix_list_23-4.start {
      counter-reset: lst-ctn-kix_list_23-4 0
    }

    .lst-kix_list_55-4 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_32-1 > li {
      counter-increment: lst-ctn-kix_list_32-1
    }

    .lst-kix_list_54-8 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) "." counter(lst-ctn-kix_list_54-2, decimal) "." counter(lst-ctn-kix_list_54-3, decimal) "." counter(lst-ctn-kix_list_54-4, decimal) "." counter(lst-ctn-kix_list_54-5, decimal) "." counter(lst-ctn-kix_list_54-6, decimal) "." counter(lst-ctn-kix_list_54-7, decimal) "." counter(lst-ctn-kix_list_54-8, decimal) " "
    }

    .lst-kix_list_13-2 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) "." counter(lst-ctn-kix_list_13-2, decimal) ". "
    }

    ol.lst-kix_list_5-2.start {
      counter-reset: lst-ctn-kix_list_5-2 0
    }

    .lst-kix_list_2-1 > li {
      counter-increment: lst-ctn-kix_list_2-1
    }

    ol.lst-kix_list_33-5.start {
      counter-reset: lst-ctn-kix_list_33-5 0
    }

    ol.lst-kix_list_9-0.start {
      counter-reset: lst-ctn-kix_list_9-0 0
    }

    ol.lst-kix_list_30-0.start {
      counter-reset: lst-ctn-kix_list_30-0 0
    }

    ol.lst-kix_list_13-4.start {
      counter-reset: lst-ctn-kix_list_13-4 0
    }

    .lst-kix_list_30-6 > li {
      counter-increment: lst-ctn-kix_list_30-6
    }

    ol.lst-kix_list_40-3.start {
      counter-reset: lst-ctn-kix_list_40-3 0
    }

    .lst-kix_list_42-8 > li:before {
      content: "" counter(lst-ctn-kix_list_42-8, decimal) " "
    }

    ol.lst-kix_list_50-1.start {
      counter-reset: lst-ctn-kix_list_50-1 0
    }

    .lst-kix_list_13-0 > li {
      counter-increment: lst-ctn-kix_list_13-0
    }

    ol.lst-kix_list_43-8.start {
      counter-reset: lst-ctn-kix_list_43-8 0
    }

    ol.lst-kix_list_20-2.start {
      counter-reset: lst-ctn-kix_list_20-2 0
    }

    .lst-kix_list_42-0 > li:before {
      content: "" counter(lst-ctn-kix_list_42-0, decimal) " "
    }

    ol.lst-kix_list_23-2.start {
      counter-reset: lst-ctn-kix_list_23-2 0
    }

    .lst-kix_list_54-0 > li {
      counter-increment: lst-ctn-kix_list_54-0
    }

    .lst-kix_list_42-2 > li:before {
      content: "" counter(lst-ctn-kix_list_42-2, decimal) " "
    }

    .lst-kix_list_43-0 > li {
      counter-increment: lst-ctn-kix_list_43-0
    }

    ol.lst-kix_list_50-6.start {
      counter-reset: lst-ctn-kix_list_50-6 0
    }

    .lst-kix_list_35-7 > li {
      counter-increment: lst-ctn-kix_list_35-7
    }

    .lst-kix_list_29-8 > li {
      counter-increment: lst-ctn-kix_list_29-8
    }

    .lst-kix_list_42-6 > li:before {
      content: "" counter(lst-ctn-kix_list_42-6, decimal) " "
    }

    .lst-kix_list_24-7 > li {
      counter-increment: lst-ctn-kix_list_24-7
    }

    ol.lst-kix_list_30-5.start {
      counter-reset: lst-ctn-kix_list_30-5 0
    }

    .lst-kix_list_42-4 > li:before {
      content: "" counter(lst-ctn-kix_list_42-4, decimal) " "
    }

    .lst-kix_list_18-8 > li {
      counter-increment: lst-ctn-kix_list_18-8
    }

    ol.lst-kix_list_26-7.start {
      counter-reset: lst-ctn-kix_list_26-7 0
    }

    ol.lst-kix_list_10-4.start {
      counter-reset: lst-ctn-kix_list_10-4 0
    }

    ol.lst-kix_list_36-5.start {
      counter-reset: lst-ctn-kix_list_36-5 0
    }

    .lst-kix_list_24-7 > li:before {
      content: "" counter(lst-ctn-kix_list_24-7, lower-latin) ". "
    }

    .lst-kix_list_1-4 > li {
      counter-increment: lst-ctn-kix_list_1-4
    }

    ol.lst-kix_list_1-6.start {
      counter-reset: lst-ctn-kix_list_1-6 0
    }

    ol.lst-kix_list_9-5.start {
      counter-reset: lst-ctn-kix_list_9-5 0
    }

    ol.lst-kix_list_40-8.start {
      counter-reset: lst-ctn-kix_list_40-8 0
    }

    ol.lst-kix_list_20-7.start {
      counter-reset: lst-ctn-kix_list_20-7 0
    }

    .lst-kix_list_24-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_24-3, lower-roman) ") "
    }

    .lst-kix_list_24-5 > li:before {
      content: "" counter(lst-ctn-kix_list_24-5, lower-latin) ") "
    }

    ol.lst-kix_list_43-3.start {
      counter-reset: lst-ctn-kix_list_43-3 0
    }

    ol.lst-kix_list_16-4.start {
      counter-reset: lst-ctn-kix_list_16-4 0
    }

    .lst-kix_list_6-5 > li {
      counter-increment: lst-ctn-kix_list_6-5
    }

    .lst-kix_list_23-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_23-3, decimal) ") "
    }

    .lst-kix_list_23-7 > li:before {
      content: "" counter(lst-ctn-kix_list_23-7, lower-latin) ". "
    }

    ol.lst-kix_list_23-7.start {
      counter-reset: lst-ctn-kix_list_23-7 0
    }

    .lst-kix_list_23-1 > li:before {
      content: "" counter(lst-ctn-kix_list_23-1, lower-latin) ") "
    }

    .lst-kix_list_24-1 > li:before {
      content: "Section " counter(lst-ctn-kix_list_24-0, upper-roman) "." counter(lst-ctn-kix_list_24-1, decimal) " "
    }

    .lst-kix_list_2-8 > li {
      counter-increment: lst-ctn-kix_list_2-8
    }

    ol.lst-kix_list_26-2.start {
      counter-reset: lst-ctn-kix_list_26-2 0
    }

    .lst-kix_list_23-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_23-5, lower-roman) ") "
    }

    ol.lst-kix_list_36-0.start {
      counter-reset: lst-ctn-kix_list_36-0 0
    }

    ol.lst-kix_list_4-6.start {
      counter-reset: lst-ctn-kix_list_4-6 0
    }

    ol.lst-kix_list_39-5.start {
      counter-reset: lst-ctn-kix_list_39-5 0
    }

    ol.lst-kix_list_3-0.start {
      counter-reset: lst-ctn-kix_list_3-0 0
    }

    ol.lst-kix_list_29-2.start {
      counter-reset: lst-ctn-kix_list_29-2 0
    }

    .lst-kix_list_25-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_25-5, lower-roman) ") "
    }

    .lst-kix_list_25-7 > li:before {
      content: "" counter(lst-ctn-kix_list_25-7, lower-latin) ". "
    }

    .lst-kix_list_3-5 > li {
      counter-increment: lst-ctn-kix_list_3-5
    }

    ol.lst-kix_list_1-1.start {
      counter-reset: lst-ctn-kix_list_1-1 0
    }

    ol.lst-kix_list_18-3.start {
      counter-reset: lst-ctn-kix_list_18-3 0
    }

    ol.lst-kix_list_33-0.start {
      counter-reset: lst-ctn-kix_list_33-0 0
    }

    ol.lst-kix_list_45-2.start {
      counter-reset: lst-ctn-kix_list_45-2 0
    }

    .lst-kix_list_9-3 > li {
      counter-increment: lst-ctn-kix_list_9-3
    }

    ol.lst-kix_list_15-3.start {
      counter-reset: lst-ctn-kix_list_15-3 0
    }

    ol.lst-kix_list_6-0.start {
      counter-reset: lst-ctn-kix_list_6-0 0
    }

    .lst-kix_list_4-2 > li {
      counter-increment: lst-ctn-kix_list_4-2
    }

    .lst-kix_list_5-1 > li {
      counter-increment: lst-ctn-kix_list_5-1
    }

    ol.lst-kix_list_21-3.start {
      counter-reset: lst-ctn-kix_list_21-3 0
    }

    .lst-kix_list_38-5 > li {
      counter-increment: lst-ctn-kix_list_38-5
    }

    ol.lst-kix_list_25-6.start {
      counter-reset: lst-ctn-kix_list_25-6 0
    }

    ol.lst-kix_list_51-7.start {
      counter-reset: lst-ctn-kix_list_51-7 0
    }

    ol.lst-kix_list_32-4.start {
      counter-reset: lst-ctn-kix_list_32-4 0
    }

    .lst-kix_list_26-5 > li:before {
      content: "" counter(lst-ctn-kix_list_26-5, lower-latin) ") "
    }

    .lst-kix_list_10-0 > li {
      counter-increment: lst-ctn-kix_list_10-0
    }

    .lst-kix_list_26-1 > li:before {
      content: "Section " counter(lst-ctn-kix_list_26-0, upper-roman) "." counter(lst-ctn-kix_list_26-1, decimal) " "
    }

    .lst-kix_list_21-1 > li:before {
      content: "" counter(lst-ctn-kix_list_21-1, decimal) ". "
    }

    ol.lst-kix_list_15-8.start {
      counter-reset: lst-ctn-kix_list_15-8 0
    }

    .lst-kix_list_34-3 > li {
      counter-increment: lst-ctn-kix_list_34-3
    }

    .lst-kix_list_45-3 > li {
      counter-increment: lst-ctn-kix_list_45-3
    }

    .lst-kix_list_10-2 > li {
      counter-increment: lst-ctn-kix_list_10-2
    }

    ol.lst-kix_list_45-7.start {
      counter-reset: lst-ctn-kix_list_45-7 0
    }

    .lst-kix_list_21-5 > li:before {
      content: "" counter(lst-ctn-kix_list_21-5, decimal) ". "
    }

    .lst-kix_list_45-2 > li:before {
      content: "" counter(lst-ctn-kix_list_45-2, decimal) " "
    }

    ol.lst-kix_list_27-8.start {
      counter-reset: lst-ctn-kix_list_27-8 0
    }

    .lst-kix_list_21-0 > li {
      counter-increment: lst-ctn-kix_list_21-0
    }

    .lst-kix_list_25-1 > li:before {
      content: "" counter(lst-ctn-kix_list_25-1, lower-latin) ") "
    }

    ol.lst-kix_list_51-2.start {
      counter-reset: lst-ctn-kix_list_51-2 0
    }

    ol.lst-kix_list_31-1.start {
      counter-reset: lst-ctn-kix_list_31-1 0
    }

    .lst-kix_list_45-6 > li:before {
      content: "" counter(lst-ctn-kix_list_45-6, decimal) " "
    }

    .lst-kix_list_44-6 > li:before {
      content: "" counter(lst-ctn-kix_list_44-6, decimal) " "
    }

    .lst-kix_list_44-2 > li:before {
      content: "" counter(lst-ctn-kix_list_44-2, decimal) " "
    }

    .lst-kix_list_39-1 > li:before {
      content: "" counter(lst-ctn-kix_list_39-0, decimal) "." counter(lst-ctn-kix_list_39-1, decimal) ". "
    }

    .lst-kix_list_16-7 > li {
      counter-increment: lst-ctn-kix_list_16-7
    }

    .lst-kix_list_27-7 > li {
      counter-increment: lst-ctn-kix_list_27-7
    }

    .lst-kix_list_45-1 > li {
      counter-increment: lst-ctn-kix_list_45-1
    }

    .lst-kix_list_34-1 > li {
      counter-increment: lst-ctn-kix_list_34-1
    }

    .lst-kix_list_38-7 > li {
      counter-increment: lst-ctn-kix_list_38-7
    }

    .lst-kix_list_16-5 > li {
      counter-increment: lst-ctn-kix_list_16-5
    }

    .lst-kix_list_3-7 > li {
      counter-increment: lst-ctn-kix_list_3-7
    }

    .lst-kix_list_50-5 > li {
      counter-increment: lst-ctn-kix_list_50-5
    }

    .lst-kix_list_22-5 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) "." counter(lst-ctn-kix_list_22-2, decimal) "." counter(lst-ctn-kix_list_22-3, decimal) "." counter(lst-ctn-kix_list_22-4, decimal) "." counter(lst-ctn-kix_list_22-5, decimal) ". "
    }

    .lst-kix_list_43-2 > li:before {
      content: "" counter(lst-ctn-kix_list_43-2, decimal) " "
    }

    .lst-kix_list_43-6 > li:before {
      content: "" counter(lst-ctn-kix_list_43-6, decimal) " "
    }

    .lst-kix_list_21-2 > li {
      counter-increment: lst-ctn-kix_list_21-2
    }

    .lst-kix_list_45-8 > li {
      counter-increment: lst-ctn-kix_list_45-8
    }

    .lst-kix_list_22-1 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) ". "
    }

    .lst-kix_list_23-3 > li {
      counter-increment: lst-ctn-kix_list_23-3
    }

    .lst-kix_list_9-5 > li {
      counter-increment: lst-ctn-kix_list_9-5
    }

    .lst-kix_list_5-8 > li {
      counter-increment: lst-ctn-kix_list_5-8
    }

    .lst-kix_list_27-0 > li {
      counter-increment: lst-ctn-kix_list_27-0
    }

    .lst-kix_list_41-6 > li {
      counter-increment: lst-ctn-kix_list_41-6
    }

    ol.lst-kix_list_19-4.start {
      counter-reset: lst-ctn-kix_list_19-4 0
    }

    .lst-kix_list_38-0 > li {
      counter-increment: lst-ctn-kix_list_38-0
    }

    .lst-kix_list_32-0 > li {
      counter-increment: lst-ctn-kix_list_32-0
    }

    ol.lst-kix_list_2-2.start {
      counter-reset: lst-ctn-kix_list_2-2 0
    }

    .lst-kix_list_50-3 > li {
      counter-increment: lst-ctn-kix_list_50-3
    }

    ol.lst-kix_list_38-4.start {
      counter-reset: lst-ctn-kix_list_38-4 0
    }

    ol.lst-kix_list_25-1.start {
      counter-reset: lst-ctn-kix_list_25-1 0
    }

    ol.lst-kix_list_21-8.start {
      counter-reset: lst-ctn-kix_list_21-8 0
    }

    .lst-kix_list_34-8 > li {
      counter-increment: lst-ctn-kix_list_34-8
    }

    .lst-kix_list_27-5 > li {
      counter-increment: lst-ctn-kix_list_27-5
    }

    .lst-kix_list_40-6 > li:before {
      content: "" counter(lst-ctn-kix_list_40-6, decimal) " "
    }

    ol.lst-kix_list_44-4.start {
      counter-reset: lst-ctn-kix_list_44-4 0
    }

    .lst-kix_list_41-6 > li:before {
      content: "" counter(lst-ctn-kix_list_41-6, decimal) " "
    }

    .lst-kix_list_20-5 > li:before {
      content: "" counter(lst-ctn-kix_list_20-5, decimal) ". "
    }

    .lst-kix_list_41-2 > li:before {
      content: "" counter(lst-ctn-kix_list_41-2, decimal) " "
    }

    .lst-kix_list_20-1 > li:before {
      content: "" counter(lst-ctn-kix_list_20-1, decimal) ". "
    }

    .lst-kix_list_12-3 > li {
      counter-increment: lst-ctn-kix_list_12-3
    }

    ol.lst-kix_list_3-5.start {
      counter-reset: lst-ctn-kix_list_3-5 0
    }

    .lst-kix_list_21-7 > li {
      counter-increment: lst-ctn-kix_list_21-7
    }

    .lst-kix_list_14-4 > li {
      counter-increment: lst-ctn-kix_list_14-4
    }

    .lst-kix_list_10-7 > li {
      counter-increment: lst-ctn-kix_list_10-7
    }

    .lst-kix_list_54-7 > li {
      counter-increment: lst-ctn-kix_list_54-7
    }

    .lst-kix_list_25-4 > li {
      counter-increment: lst-ctn-kix_list_25-4
    }

    .lst-kix_list_18-1 > li {
      counter-increment: lst-ctn-kix_list_18-1
    }

    ol.lst-kix_list_14-5.start {
      counter-reset: lst-ctn-kix_list_14-5 0
    }

    .lst-kix_list_32-7 > li {
      counter-increment: lst-ctn-kix_list_32-7
    }

    .lst-kix_list_40-2 > li:before {
      content: "" counter(lst-ctn-kix_list_40-2, decimal) " "
    }

    ol.lst-kix_list_37-1.start {
      counter-reset: lst-ctn-kix_list_37-1 0
    }

    .lst-kix_list_36-4 > li {
      counter-increment: lst-ctn-kix_list_36-4
    }

    .lst-kix_list_29-1 > li {
      counter-increment: lst-ctn-kix_list_29-1
    }

    .lst-kix_list_43-7 > li {
      counter-increment: lst-ctn-kix_list_43-7
    }

    ol.lst-kix_list_32-7.start {
      counter-reset: lst-ctn-kix_list_32-7 0
    }

    .lst-kix_list_43-2 > li {
      counter-increment: lst-ctn-kix_list_43-2
    }

    ol.lst-kix_list_41-1.start {
      counter-reset: lst-ctn-kix_list_41-1 0
    }

    .lst-kix_list_19-4 > li:before {
      content: "" counter(lst-ctn-kix_list_19-4, decimal) ". "
    }

    ol.lst-kix_list_37-8.start {
      counter-reset: lst-ctn-kix_list_37-8 0
    }

    .lst-kix_list_31-5 > li {
      counter-increment: lst-ctn-kix_list_31-5
    }

    ol.lst-kix_list_24-5.start {
      counter-reset: lst-ctn-kix_list_24-5 0
    }

    .lst-kix_list_47-6 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_47-8 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_27-5.start {
      counter-reset: lst-ctn-kix_list_27-5 0
    }

    ol.lst-kix_list_22-4.start {
      counter-reset: lst-ctn-kix_list_22-4 0
    }

    .lst-kix_list_42-1 > li {
      counter-increment: lst-ctn-kix_list_42-1
    }

    ol.lst-kix_list_11-2.start {
      counter-reset: lst-ctn-kix_list_11-2 0
    }

    .lst-kix_list_19-6 > li:before {
      content: "" counter(lst-ctn-kix_list_19-6, decimal) ". "
    }

    .lst-kix_list_47-0 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_17-2 > li {
      counter-increment: lst-ctn-kix_list_17-2
    }

    ol.lst-kix_list_41-6.start {
      counter-reset: lst-ctn-kix_list_41-6 0
    }

    .lst-kix_list_42-5 > li {
      counter-increment: lst-ctn-kix_list_42-5
    }

    .lst-kix_list_32-2 > li {
      counter-increment: lst-ctn-kix_list_32-2
    }

    .lst-kix_list_36-2 > li {
      counter-increment: lst-ctn-kix_list_36-2
    }

    .lst-kix_list_20-5 > li {
      counter-increment: lst-ctn-kix_list_20-5
    }

    ol.lst-kix_list_37-3.start {
      counter-reset: lst-ctn-kix_list_37-3 0
    }

    .lst-kix_list_37-1 > li:before {
      content: "" counter(lst-ctn-kix_list_37-0, decimal) "." counter(lst-ctn-kix_list_37-1, decimal) ". "
    }

    .lst-kix_list_46-2 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_44-1.start {
      counter-reset: lst-ctn-kix_list_44-1 0
    }

    .lst-kix_list_54-2 > li {
      counter-increment: lst-ctn-kix_list_54-2
    }

    ol.lst-kix_list_24-0.start {
      counter-reset: lst-ctn-kix_list_24-0 0
    }

    .lst-kix_list_18-3 > li {
      counter-increment: lst-ctn-kix_list_18-3
    }

    .lst-kix_list_37-3 > li:before {
      content: "" counter(lst-ctn-kix_list_37-0, decimal) "." counter(lst-ctn-kix_list_37-1, decimal) "." counter(lst-ctn-kix_list_37-2, decimal) "." counter(lst-ctn-kix_list_37-3, decimal) " "
    }

    ol.lst-kix_list_35-7.start {
      counter-reset: lst-ctn-kix_list_35-7 0
    }

    .lst-kix_list_37-3 > li {
      counter-increment: lst-ctn-kix_list_37-3
    }

    .lst-kix_list_18-0 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_38-7 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) "." counter(lst-ctn-kix_list_38-2, decimal) "." counter(lst-ctn-kix_list_38-3, decimal) "." counter(lst-ctn-kix_list_38-4, decimal) "." counter(lst-ctn-kix_list_38-5, decimal) "." counter(lst-ctn-kix_list_38-6, decimal) "." counter(lst-ctn-kix_list_38-7, decimal) " "
    }

    .lst-kix_list_3-0 > li {
      counter-increment: lst-ctn-kix_list_3-0
    }

    .lst-kix_list_18-2 > li:before {
      content: "" counter(lst-ctn-kix_list_18-2, decimal) ". "
    }

    ol.lst-kix_list_27-0.start {
      counter-reset: lst-ctn-kix_list_27-0 0
    }

    ol.lst-kix_list_54-4.start {
      counter-reset: lst-ctn-kix_list_54-4 0
    }

    ol.lst-kix_list_11-7.start {
      counter-reset: lst-ctn-kix_list_11-7 0
    }

    ol.lst-kix_list_14-2.start {
      counter-reset: lst-ctn-kix_list_14-2 0
    }

    .lst-kix_list_41-4 > li {
      counter-increment: lst-ctn-kix_list_41-4
    }

    .lst-kix_list_30-4 > li {
      counter-increment: lst-ctn-kix_list_30-4
    }

    .lst-kix_list_38-5 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) "." counter(lst-ctn-kix_list_38-2, decimal) "." counter(lst-ctn-kix_list_38-3, decimal) "." counter(lst-ctn-kix_list_38-4, decimal) "." counter(lst-ctn-kix_list_38-5, decimal) " "
    }

    .lst-kix_list_25-2 > li {
      counter-increment: lst-ctn-kix_list_25-2
    }

    .lst-kix_list_27-1 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) ". "
    }

    .lst-kix_list_48-2 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_27-3 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) "." counter(lst-ctn-kix_list_27-2, decimal) "." counter(lst-ctn-kix_list_27-3, decimal) ". "
    }

    .lst-kix_list_48-4 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_18-8 > li:before {
      content: "" counter(lst-ctn-kix_list_18-8, decimal) ". "
    }

    ol.lst-kix_list_27-3.start {
      counter-reset: lst-ctn-kix_list_27-3 0
    }

    .lst-kix_list_17-6 > li {
      counter-increment: lst-ctn-kix_list_17-6
    }

    ol.lst-kix_list_4-3.start {
      counter-reset: lst-ctn-kix_list_4-3 0
    }

    .lst-kix_list_10-7 > li:before {
      content: "" counter(lst-ctn-kix_list_10-7, decimal) " "
    }

    .lst-kix_list_20-1 > li {
      counter-increment: lst-ctn-kix_list_20-1
    }

    .lst-kix_list_10-5 > li:before {
      content: "" counter(lst-ctn-kix_list_10-5, decimal) " "
    }

    ol.lst-kix_list_18-6.start {
      counter-reset: lst-ctn-kix_list_18-6 0
    }

    .lst-kix_list_29-3 > li {
      counter-increment: lst-ctn-kix_list_29-3
    }

    ol.lst-kix_list_54-2.start {
      counter-reset: lst-ctn-kix_list_54-2 0
    }

    .lst-kix_list_9-2 > li:before {
      content: "" counter(lst-ctn-kix_list_9-0, decimal) "." counter(lst-ctn-kix_list_9-1, decimal) "." counter(lst-ctn-kix_list_9-2, decimal) " "
    }

    .lst-kix_list_46-4 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_37-6.start {
      counter-reset: lst-ctn-kix_list_37-6 0
    }

    ol.lst-kix_list_14-0.start {
      counter-reset: lst-ctn-kix_list_14-0 0
    }

    .lst-kix_list_12-5 > li {
      counter-increment: lst-ctn-kix_list_12-5
    }

    .lst-kix_list_31-1 > li {
      counter-increment: lst-ctn-kix_list_31-1
    }

    .lst-kix_list_9-0 > li:before {
      content: "" counter(lst-ctn-kix_list_9-0, decimal) " "
    }

    ol.lst-kix_list_24-3.start {
      counter-reset: lst-ctn-kix_list_24-3 0
    }

    .lst-kix_list_23-5 > li {
      counter-increment: lst-ctn-kix_list_23-5
    }

    .lst-kix_list_11-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_11-3, lower-latin) ") "
    }

    .lst-kix_list_6-3 > li {
      counter-increment: lst-ctn-kix_list_6-3
    }

    ol.lst-kix_list_1-3.start {
      counter-reset: lst-ctn-kix_list_1-3 0
    }

    .lst-kix_list_29-1 > li:before {
      content: "Section " counter(lst-ctn-kix_list_29-0, upper-roman) "." counter(lst-ctn-kix_list_29-1, decimal) " "
    }

    .lst-kix_list_29-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_29-3, lower-roman) ") "
    }

    .lst-kix_list_49-8 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_9-8 > li:before {
      content: "" counter(lst-ctn-kix_list_9-8, decimal) " "
    }

    .lst-kix_list_28-6 > li {
      counter-increment: lst-ctn-kix_list_28-6
    }

    .lst-kix_list_28-7 > li:before {
      content: "" counter(lst-ctn-kix_list_28-7, lower-latin) ". "
    }

    .lst-kix_list_1-7 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) "." counter(lst-ctn-kix_list_1-2, decimal) "." counter(lst-ctn-kix_list_1-3, decimal) "." counter(lst-ctn-kix_list_1-4, decimal) "." counter(lst-ctn-kix_list_1-5, decimal) "." counter(lst-ctn-kix_list_1-6, decimal) "." counter(lst-ctn-kix_list_1-7, decimal) " "
    }

    .lst-kix_list_50-1 > li {
      counter-increment: lst-ctn-kix_list_50-1
    }

    .lst-kix_list_49-6 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_1-5 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) "." counter(lst-ctn-kix_list_1-2, decimal) "." counter(lst-ctn-kix_list_1-3, decimal) "." counter(lst-ctn-kix_list_1-4, decimal) "." counter(lst-ctn-kix_list_1-5, decimal) " "
    }

    .lst-kix_list_28-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_28-5, lower-roman) ") "
    }

    .lst-kix_list_5-6 > li {
      counter-increment: lst-ctn-kix_list_5-6
    }

    ol.lst-kix_list_22-6.start {
      counter-reset: lst-ctn-kix_list_22-6 0
    }

    .lst-kix_list_2-1 > li:before {
      content: "" counter(lst-ctn-kix_list_2-1, decimal) ". "
    }

    .lst-kix_list_49-0 > li:before {
      content: "\0025a0  "
    }

    .lst-kix_list_2-3 > li:before {
      content: "" counter(lst-ctn-kix_list_2-3, decimal) ". "
    }

    .lst-kix_list_11-8 > li {
      counter-increment: lst-ctn-kix_list_11-8
    }

    .lst-kix_list_35-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_35-3, lower-latin) ") "
    }

    ol.lst-kix_list_24-8.start {
      counter-reset: lst-ctn-kix_list_24-8 0
    }

    .lst-kix_list_30-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_30-4, lower-roman) ") "
    }

    .lst-kix_list_20-8 > li {
      counter-increment: lst-ctn-kix_list_20-8
    }

    .lst-kix_list_9-1 > li {
      counter-increment: lst-ctn-kix_list_9-1
    }

    .lst-kix_list_26-7 > li:before {
      content: "" counter(lst-ctn-kix_list_26-7, lower-latin) ". "
    }

    ol.lst-kix_list_1-8.start {
      counter-reset: lst-ctn-kix_list_1-8 0
    }

    .lst-kix_list_6-0 > li {
      counter-increment: lst-ctn-kix_list_6-0
    }

    .lst-kix_list_3-5 > li:before {
      content: "" counter(lst-ctn-kix_list_3-5, decimal) ". "
    }

    .lst-kix_list_40-2 > li {
      counter-increment: lst-ctn-kix_list_40-2
    }

    ol.lst-kix_list_30-8.start {
      counter-reset: lst-ctn-kix_list_30-8 0
    }

    ol.lst-kix_list_11-5.start {
      counter-reset: lst-ctn-kix_list_11-5 0
    }

    .lst-kix_list_8-6 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_11-1 > li {
      counter-increment: lst-ctn-kix_list_11-1
    }

    ol.lst-kix_list_16-6.start {
      counter-reset: lst-ctn-kix_list_16-6 0
    }

    .lst-kix_list_44-4 > li {
      counter-increment: lst-ctn-kix_list_44-4
    }

    ol.lst-kix_list_41-4.start {
      counter-reset: lst-ctn-kix_list_41-4 0
    }

    ol.lst-kix_list_22-1.start {
      counter-reset: lst-ctn-kix_list_22-1 0
    }

    .lst-kix_list_53-0 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_33-4 > li {
      counter-increment: lst-ctn-kix_list_33-4
    }

    .lst-kix_list_21-3 > li:before {
      content: "" counter(lst-ctn-kix_list_21-3, decimal) ". "
    }

    ol.lst-kix_list_30-2.start {
      counter-reset: lst-ctn-kix_list_30-2 0
    }

    .lst-kix_list_4-4 > li {
      counter-increment: lst-ctn-kix_list_4-4
    }

    ol.lst-kix_list_29-4.start {
      counter-reset: lst-ctn-kix_list_29-4 0
    }

    .lst-kix_list_45-8 > li:before {
      content: "" counter(lst-ctn-kix_list_45-8, decimal) " "
    }

    .lst-kix_list_31-8 > li {
      counter-increment: lst-ctn-kix_list_31-8
    }

    .lst-kix_list_25-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_25-3, decimal) ") "
    }

    ol.lst-kix_list_4-1.start {
      counter-reset: lst-ctn-kix_list_4-1 0
    }

    .lst-kix_list_39-2 > li {
      counter-increment: lst-ctn-kix_list_39-2
    }

    .lst-kix_list_16-2 > li:before {
      content: "" counter(lst-ctn-kix_list_16-2, decimal) ". "
    }

    .lst-kix_list_26-6 > li {
      counter-increment: lst-ctn-kix_list_26-6
    }

    .lst-kix_list_15-3 > li {
      counter-increment: lst-ctn-kix_list_15-3
    }

    .lst-kix_list_44-4 > li:before {
      content: "" counter(lst-ctn-kix_list_44-4, decimal) " "
    }

    ol.lst-kix_list_30-3.start {
      counter-reset: lst-ctn-kix_list_30-3 0
    }

    .lst-kix_list_37-6 > li {
      counter-increment: lst-ctn-kix_list_37-6
    }

    ol.lst-kix_list_11-0.start {
      counter-reset: lst-ctn-kix_list_11-0 0
    }

    ol.lst-kix_list_18-8.start {
      counter-reset: lst-ctn-kix_list_18-8 0
    }

    .lst-kix_list_3-3 > li {
      counter-increment: lst-ctn-kix_list_3-3
    }

    .lst-kix_list_39-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_39-3, lower-latin) ") "
    }

    ol.lst-kix_list_6-3.start {
      counter-reset: lst-ctn-kix_list_6-3 0
    }

    .lst-kix_list_45-0 > li:before {
      content: "" counter(lst-ctn-kix_list_45-0, decimal) " "
    }

    .lst-kix_list_17-6 > li:before {
      content: "" counter(lst-ctn-kix_list_17-6, decimal) ". "
    }

    ol.lst-kix_list_16-2.start {
      counter-reset: lst-ctn-kix_list_16-2 0
    }

    .lst-kix_list_42-8 > li {
      counter-increment: lst-ctn-kix_list_42-8
    }

    .lst-kix_list_43-0 > li:before {
      content: "" counter(lst-ctn-kix_list_43-0, decimal) " "
    }

    .lst-kix_list_43-8 > li:before {
      content: "" counter(lst-ctn-kix_list_43-8, decimal) " "
    }

    .lst-kix_list_19-5 > li {
      counter-increment: lst-ctn-kix_list_19-5
    }

    .lst-kix_list_28-2 > li {
      counter-increment: lst-ctn-kix_list_28-2
    }

    .lst-kix_list_22-7 > li:before {
      content: "" counter(lst-ctn-kix_list_22-0, decimal) "." counter(lst-ctn-kix_list_22-1, decimal) "." counter(lst-ctn-kix_list_22-2, decimal) "." counter(lst-ctn-kix_list_22-3, decimal) "." counter(lst-ctn-kix_list_22-4, decimal) "." counter(lst-ctn-kix_list_22-5, decimal) "." counter(lst-ctn-kix_list_22-6, decimal) "." counter(lst-ctn-kix_list_22-7, decimal) ". "
    }

    ol.lst-kix_list_35-2.start {
      counter-reset: lst-ctn-kix_list_35-2 0
    }

    .lst-kix_list_34-7 > li:before {
      content: "" counter(lst-ctn-kix_list_34-7, decimal) " "
    }

    .lst-kix_list_55-8 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_6-7 > li {
      counter-increment: lst-ctn-kix_list_6-7
    }

    .lst-kix_list_26-3 > li {
      counter-increment: lst-ctn-kix_list_26-3
    }

    .lst-kix_list_15-6 > li:before {
      content: "" counter(lst-ctn-kix_list_15-6, lower-roman) ") "
    }

    .lst-kix_list_11-4 > li {
      counter-increment: lst-ctn-kix_list_11-4
    }

    .lst-kix_list_22-4 > li {
      counter-increment: lst-ctn-kix_list_22-4
    }

    .lst-kix_list_52-4 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_36-7 > li:before {
      content: "" counter(lst-ctn-kix_list_36-7, lower-latin) ". "
    }

    ol.lst-kix_list_6-8.start {
      counter-reset: lst-ctn-kix_list_6-8 0
    }

    .lst-kix_list_13-7 > li {
      counter-increment: lst-ctn-kix_list_13-7
    }

    .lst-kix_list_20-7 > li:before {
      content: "" counter(lst-ctn-kix_list_20-7, decimal) ". "
    }

    ol.lst-kix_list_6-5.start {
      counter-reset: lst-ctn-kix_list_6-5 0
    }

    .lst-kix_list_41-8 > li:before {
      content: "" counter(lst-ctn-kix_list_41-8, decimal) " "
    }

    ol.lst-kix_list_29-5.start {
      counter-reset: lst-ctn-kix_list_29-5 0
    }

    .lst-kix_list_51-2 > li {
      counter-increment: lst-ctn-kix_list_51-2
    }

    .lst-kix_list_41-0 > li:before {
      content: "" counter(lst-ctn-kix_list_41-0, decimal) " "
    }

    ol.lst-kix_list_6-7.start {
      counter-reset: lst-ctn-kix_list_6-7 0
    }

    .lst-kix_list_24-0 > li {
      counter-increment: lst-ctn-kix_list_24-0
    }

    .lst-kix_list_54-4 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) "." counter(lst-ctn-kix_list_54-2, decimal) "." counter(lst-ctn-kix_list_54-3, decimal) "." counter(lst-ctn-kix_list_54-4, decimal) " "
    }

    .lst-kix_list_33-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_33-3, lower-latin) ") "
    }

    .lst-kix_list_53-8 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_55-0 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_16-1.start {
      counter-reset: lst-ctn-kix_list_16-1 0
    }

    ol.lst-kix_list_35-4.start {
      counter-reset: lst-ctn-kix_list_35-4 0
    }

    ol.lst-kix_list_54-7.start {
      counter-reset: lst-ctn-kix_list_54-7 0
    }

    ol.lst-kix_list_29-7.start {
      counter-reset: lst-ctn-kix_list_29-7 0
    }

    .lst-kix_list_35-0 > li {
      counter-increment: lst-ctn-kix_list_35-0
    }

    .lst-kix_list_40-4 > li:before {
      content: "" counter(lst-ctn-kix_list_40-4, decimal) " "
    }

    ol.lst-kix_list_30-7.start {
      counter-reset: lst-ctn-kix_list_30-7 0
    }

    .lst-kix_list_51-5 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) "." counter(lst-ctn-kix_list_51-2, decimal) "." counter(lst-ctn-kix_list_51-3, decimal) "." counter(lst-ctn-kix_list_51-4, decimal) "." counter(lst-ctn-kix_list_51-5, decimal) " "
    }

    .lst-kix_list_41-2 > li {
      counter-increment: lst-ctn-kix_list_41-2
    }

    .lst-kix_list_14-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_14-3, decimal) ") "
    }

    .lst-kix_list_14-0 > li:before {
      content: "" counter(lst-ctn-kix_list_14-0, decimal) ") "
    }

    ol.lst-kix_list_18-5.start {
      counter-reset: lst-ctn-kix_list_18-5 0
    }

    .lst-kix_list_6-1 > li {
      counter-increment: lst-ctn-kix_list_6-1
    }

    .lst-kix_list_51-2 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) "." counter(lst-ctn-kix_list_51-2, decimal) " "
    }

    .lst-kix_list_51-3 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) "." counter(lst-ctn-kix_list_51-2, decimal) "." counter(lst-ctn-kix_list_51-3, decimal) " "
    }

    .lst-kix_list_14-6 > li:before {
      content: "" counter(lst-ctn-kix_list_14-6, decimal) ". "
    }

    .lst-kix_ek190fjkpfvz-8 > li:before {
      content: "\0025a0  "
    }

    .lst-kix_list_9-0 > li {
      counter-increment: lst-ctn-kix_list_9-0
    }

    ol.lst-kix_list_25-3.start {
      counter-reset: lst-ctn-kix_list_25-3 0
    }

    .lst-kix_list_51-6 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) "." counter(lst-ctn-kix_list_51-1, decimal) "." counter(lst-ctn-kix_list_51-2, decimal) "." counter(lst-ctn-kix_list_51-3, decimal) "." counter(lst-ctn-kix_list_51-4, decimal) "." counter(lst-ctn-kix_list_51-5, decimal) "." counter(lst-ctn-kix_list_51-6, decimal) " "
    }

    ol.lst-kix_list_11-8.start {
      counter-reset: lst-ctn-kix_list_11-8 0
    }

    .lst-kix_list_14-2 > li:before {
      content: "" counter(lst-ctn-kix_list_14-2, lower-roman) ") "
    }

    .lst-kix_ek190fjkpfvz-1 > li:before {
      content: "\0025cb  "
    }

    .lst-kix_ek190fjkpfvz-2 > li:before {
      content: "\0025a0  "
    }

    ol.lst-kix_list_12-0.start {
      counter-reset: lst-ctn-kix_list_12-0 0
    }

    .lst-kix_ek190fjkpfvz-5 > li:before {
      content: "\0025a0  "
    }

    ol.lst-kix_list_41-7.start {
      counter-reset: lst-ctn-kix_list_41-7 0
    }

    .lst-kix_list_32-2 > li:before {
      content: "" counter(lst-ctn-kix_list_32-0, decimal) "." counter(lst-ctn-kix_list_32-1, decimal) "." counter(lst-ctn-kix_list_32-2, decimal) " "
    }

    .lst-kix_list_32-1 > li:before {
      content: "" counter(lst-ctn-kix_list_32-1, decimal) ". "
    }

    ol.lst-kix_list_21-6.start {
      counter-reset: lst-ctn-kix_list_21-6 0
    }

    ol.lst-kix_list_3-7.start {
      counter-reset: lst-ctn-kix_list_3-7 0
    }

    .lst-kix_list_28-7 > li {
      counter-increment: lst-ctn-kix_list_28-7
    }

    .lst-kix_list_31-7 > li {
      counter-increment: lst-ctn-kix_list_31-7
    }

    .lst-kix_list_5-0 > li:before {
      content: "" counter(lst-ctn-kix_list_5-0, decimal) ". "
    }

    .lst-kix_list_14-8 > li {
      counter-increment: lst-ctn-kix_list_14-8
    }

    .lst-kix_list_5-3 > li:before {
      content: "" counter(lst-ctn-kix_list_5-3, decimal) ". "
    }

    .lst-kix_list_36-8 > li {
      counter-increment: lst-ctn-kix_list_36-8
    }

    .lst-kix_list_5-2 > li:before {
      content: "" counter(lst-ctn-kix_list_5-2, decimal) ". "
    }

    .lst-kix_list_5-6 > li:before {
      content: "" counter(lst-ctn-kix_list_5-6, decimal) ". "
    }

    ol.lst-kix_list_12-5.start {
      counter-reset: lst-ctn-kix_list_12-5 0
    }

    .lst-kix_list_50-0 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) " "
    }

    .lst-kix_list_50-3 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) "." counter(lst-ctn-kix_list_50-2, decimal) "." counter(lst-ctn-kix_list_50-3, decimal) " "
    }

    .lst-kix_list_50-4 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) "." counter(lst-ctn-kix_list_50-2, decimal) "." counter(lst-ctn-kix_list_50-3, decimal) "." counter(lst-ctn-kix_list_50-4, decimal) " "
    }

    .lst-kix_list_6-2 > li:before {
      content: "" counter(lst-ctn-kix_list_6-2, decimal) ". "
    }

    ol.lst-kix_list_35-1.start {
      counter-reset: lst-ctn-kix_list_35-1 0
    }

    ol.lst-kix_list_3-2.start {
      counter-reset: lst-ctn-kix_list_3-2 0
    }

    .lst-kix_list_50-7 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) "." counter(lst-ctn-kix_list_50-2, decimal) "." counter(lst-ctn-kix_list_50-3, decimal) "." counter(lst-ctn-kix_list_50-4, decimal) "." counter(lst-ctn-kix_list_50-5, decimal) "." counter(lst-ctn-kix_list_50-6, decimal) "." counter(lst-ctn-kix_list_50-7, decimal) " "
    }

    .lst-kix_list_6-8 > li:before {
      content: "" counter(lst-ctn-kix_list_6-8, decimal) ". "
    }

    .lst-kix_list_51-0 > li:before {
      content: "" counter(lst-ctn-kix_list_51-0, decimal) " "
    }

    .lst-kix_list_50-6 > li:before {
      content: "" counter(lst-ctn-kix_list_50-0, decimal) "." counter(lst-ctn-kix_list_50-1, decimal) "." counter(lst-ctn-kix_list_50-2, decimal) "." counter(lst-ctn-kix_list_50-3, decimal) "." counter(lst-ctn-kix_list_50-4, decimal) "." counter(lst-ctn-kix_list_50-5, decimal) "." counter(lst-ctn-kix_list_50-6, decimal) " "
    }

    .lst-kix_list_6-5 > li:before {
      content: "" counter(lst-ctn-kix_list_6-5, decimal) ". "
    }

    ol.lst-kix_list_42-4.start {
      counter-reset: lst-ctn-kix_list_42-4 0
    }

    .lst-kix_list_7-4 > li:before {
      content: "o  "
    }

    .lst-kix_list_22-2 > li {
      counter-increment: lst-ctn-kix_list_22-2
    }

    .lst-kix_list_44-8 > li {
      counter-increment: lst-ctn-kix_list_44-8
    }

    .lst-kix_list_52-1 > li:before {
      content: "o  "
    }

    .lst-kix_list_24-3 > li {
      counter-increment: lst-ctn-kix_list_24-3
    }

    .lst-kix_list_13-7 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) "." counter(lst-ctn-kix_list_13-2, decimal) "." counter(lst-ctn-kix_list_13-3, decimal) "." counter(lst-ctn-kix_list_13-4, decimal) "." counter(lst-ctn-kix_list_13-5, decimal) "." counter(lst-ctn-kix_list_13-6, decimal) "." counter(lst-ctn-kix_list_13-7, decimal) ". "
    }

    ol.lst-kix_list_34-4.start {
      counter-reset: lst-ctn-kix_list_34-4 0
    }

    .lst-kix_list_7-8 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_15-6 > li {
      counter-increment: lst-ctn-kix_list_15-6
    }

    .lst-kix_list_4-7 > li {
      counter-increment: lst-ctn-kix_list_4-7
    }

    .lst-kix_list_51-5 > li {
      counter-increment: lst-ctn-kix_list_51-5
    }

    ol.lst-kix_list_2-5.start {
      counter-reset: lst-ctn-kix_list_2-5 0
    }

    .lst-kix_list_15-5 > li:before {
      content: "" counter(lst-ctn-kix_list_15-5, lower-latin) ") "
    }

    .lst-kix_list_31-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_31-6, decimal) ") "
    }

    .lst-kix_list_31-8 > li:before {
      content: "" counter(lst-ctn-kix_list_31-8, decimal) " "
    }

    .lst-kix_list_45-6 > li {
      counter-increment: lst-ctn-kix_list_45-6
    }

    ol.lst-kix_list_26-0.start {
      counter-reset: lst-ctn-kix_list_26-0 0
    }

    .lst-kix_list_52-5 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_27-3 > li {
      counter-increment: lst-ctn-kix_list_27-3
    }

    .lst-kix_list_54-3 > li {
      counter-increment: lst-ctn-kix_list_54-3
    }

    .lst-kix_list_4-1 > li:before {
      content: "" counter(lst-ctn-kix_list_4-1, decimal) ". "
    }

    .lst-kix_list_31-2 > li:before {
      content: "" counter(lst-ctn-kix_list_31-0, decimal) "." counter(lst-ctn-kix_list_31-1, decimal) "." counter(lst-ctn-kix_list_31-2, decimal) " "
    }

    .lst-kix_list_36-0 > li {
      counter-increment: lst-ctn-kix_list_36-0
    }

    ol.lst-kix_list_33-3.start {
      counter-reset: lst-ctn-kix_list_33-3 0
    }

    .lst-kix_list_4-3 > li:before {
      content: "" counter(lst-ctn-kix_list_4-3, decimal) ". "
    }

    .lst-kix_list_1-8 > li {
      counter-increment: lst-ctn-kix_list_1-8
    }

    ol.lst-kix_list_38-1.start {
      counter-reset: lst-ctn-kix_list_38-1 0
    }

    .lst-kix_list_10-5 > li {
      counter-increment: lst-ctn-kix_list_10-5
    }

    .lst-kix_list_24-4 > li {
      counter-increment: lst-ctn-kix_list_24-4
    }

    .lst-kix_list_33-1 > li {
      counter-increment: lst-ctn-kix_list_33-1
    }

    .lst-kix_list_16-2 > li {
      counter-increment: lst-ctn-kix_list_16-2
    }

    ol.lst-kix_list_39-3.start {
      counter-reset: lst-ctn-kix_list_39-3 0
    }

    ol.lst-kix_list_16-7.start {
      counter-reset: lst-ctn-kix_list_16-7 0
    }

    .lst-kix_list_32-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_32-4, lower-roman) ") "
    }

    .lst-kix_list_26-5 > li {
      counter-increment: lst-ctn-kix_list_26-5
    }

    .lst-kix_list_19-2 > li {
      counter-increment: lst-ctn-kix_list_19-2
    }

    .lst-kix_list_33-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_33-4, lower-roman) ") "
    }

    .lst-kix_list_30-3 > li {
      counter-increment: lst-ctn-kix_list_30-3
    }

    .lst-kix_list_38-2 > li {
      counter-increment: lst-ctn-kix_list_38-2
    }

    ol.lst-kix_list_43-6.start {
      counter-reset: lst-ctn-kix_list_43-6 0
    }

    .lst-kix_list_45-5 > li {
      counter-increment: lst-ctn-kix_list_45-5
    }

    .lst-kix_list_12-1 > li:before {
      content: "(" counter(lst-ctn-kix_list_12-1, lower-latin) ") "
    }

    .lst-kix_list_33-0 > li:before {
      content: "" counter(lst-ctn-kix_list_33-0, decimal) " "
    }

    .lst-kix_list_53-7 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_32-8 > li:before {
      content: "" counter(lst-ctn-kix_list_32-8, decimal) " "
    }

    .lst-kix_list_23-6 > li {
      counter-increment: lst-ctn-kix_list_23-6
    }

    ol.lst-kix_list_13-6.start {
      counter-reset: lst-ctn-kix_list_13-6 0
    }

    ol.lst-kix_list_25-8.start {
      counter-reset: lst-ctn-kix_list_25-8 0
    }

    .lst-kix_list_39-0 > li {
      counter-increment: lst-ctn-kix_list_39-0
    }

    .lst-kix_list_34-0 > li:before {
      content: "" counter(lst-ctn-kix_list_34-0, decimal) " "
    }

    .lst-kix_list_21-4 > li {
      counter-increment: lst-ctn-kix_list_21-4
    }

    .lst-kix_list_39-6 > li {
      counter-increment: lst-ctn-kix_list_39-6
    }

    .lst-kix_list_13-3 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) "." counter(lst-ctn-kix_list_13-2, decimal) "." counter(lst-ctn-kix_list_13-3, decimal) ". "
    }

    .lst-kix_list_42-6 > li {
      counter-increment: lst-ctn-kix_list_42-6
    }

    ol.lst-kix_list_43-5.start {
      counter-reset: lst-ctn-kix_list_43-5 0
    }

    .lst-kix_list_34-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_34-6, decimal) ") "
    }

    .lst-kix_list_12-5 > li:before {
      content: "" counter(lst-ctn-kix_list_12-5, decimal) " "
    }

    ol.lst-kix_list_13-7.start {
      counter-reset: lst-ctn-kix_list_13-7 0
    }

    .lst-kix_list_43-4 > li {
      counter-increment: lst-ctn-kix_list_43-4
    }

    .lst-kix_list_18-4 > li {
      counter-increment: lst-ctn-kix_list_18-4
    }

    .lst-kix_list_42-0 > li {
      counter-increment: lst-ctn-kix_list_42-0
    }

    .lst-kix_list_12-7 > li:before {
      content: "" counter(lst-ctn-kix_list_12-7, decimal) " "
    }

    .lst-kix_list_50-7 > li {
      counter-increment: lst-ctn-kix_list_50-7
    }

    ol.lst-kix_list_51-0.start {
      counter-reset: lst-ctn-kix_list_51-0 0
    }

    ol.lst-kix_list_21-1.start {
      counter-reset: lst-ctn-kix_list_21-1 0
    }

    ol.lst-kix_list_50-8.start {
      counter-reset: lst-ctn-kix_list_50-8 0
    }

    .lst-kix_list_25-1 > li {
      counter-increment: lst-ctn-kix_list_25-1
    }

    .lst-kix_list_13-1 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) ". "
    }

    .lst-kix_list_32-5 > li {
      counter-increment: lst-ctn-kix_list_32-5
    }

    .lst-kix_list_22-8 > li {
      counter-increment: lst-ctn-kix_list_22-8
    }

    .lst-kix_list_35-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_35-5, upper-latin) ") "
    }

    ol.lst-kix_list_33-2.start {
      counter-reset: lst-ctn-kix_list_33-2 0
    }

    ol.lst-kix_list_2-6.start {
      counter-reset: lst-ctn-kix_list_2-6 0
    }

    ol.lst-kix_list_20-5.start {
      counter-reset: lst-ctn-kix_list_20-5 0
    }

    ol.lst-kix_list_13-1.start {
      counter-reset: lst-ctn-kix_list_13-1 0
    }

    .lst-kix_list_3-4 > li:before {
      content: "" counter(lst-ctn-kix_list_3-4, decimal) ". "
    }

    .lst-kix_list_8-0 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_30-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_30-6, decimal) ") "
    }

    .lst-kix_list_8-7 > li:before {
      content: "o  "
    }

    .lst-kix_list_3-8 > li:before {
      content: "" counter(lst-ctn-kix_list_3-8, decimal) ". "
    }

    .lst-kix_list_8-3 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_43-0.start {
      counter-reset: lst-ctn-kix_list_43-0 0
    }

    .lst-kix_list_3-7 > li:before {
      content: "" counter(lst-ctn-kix_list_3-7, decimal) ". "
    }

    .lst-kix_list_8-4 > li:before {
      content: "o  "
    }

    .lst-kix_list_19-1 > li {
      counter-increment: lst-ctn-kix_list_19-1
    }

    ol.lst-kix_list_26-4.start {
      counter-reset: lst-ctn-kix_list_26-4 0
    }

    .lst-kix_list_35-8 > li:before {
      content: "" counter(lst-ctn-kix_list_35-8, decimal) " "
    }

    ol.lst-kix_list_50-3.start {
      counter-reset: lst-ctn-kix_list_50-3 0
    }

    .lst-kix_list_37-4 > li {
      counter-increment: lst-ctn-kix_list_37-4
    }

    .lst-kix_list_16-8 > li:before {
      content: "" counter(lst-ctn-kix_list_16-8, decimal) ". "
    }

    .lst-kix_list_16-7 > li:before {
      content: "" counter(lst-ctn-kix_list_16-7, decimal) ". "
    }

    .lst-kix_list_17-8 > li {
      counter-increment: lst-ctn-kix_list_17-8
    }

    .lst-kix_list_50-8 > li {
      counter-increment: lst-ctn-kix_list_50-8
    }

    .lst-kix_list_4-7 > li:before {
      content: "" counter(lst-ctn-kix_list_4-7, decimal) ". "
    }

    .lst-kix_list_17-0 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_39-8 > li {
      counter-increment: lst-ctn-kix_list_39-8
    }

    .lst-kix_list_16-4 > li:before {
      content: "" counter(lst-ctn-kix_list_16-4, decimal) ". "
    }

    ol.lst-kix_list_3-3.start {
      counter-reset: lst-ctn-kix_list_3-3 0
    }

    .lst-kix_list_11-3 > li {
      counter-increment: lst-ctn-kix_list_11-3
    }

    ol.lst-kix_list_50-4.start {
      counter-reset: lst-ctn-kix_list_50-4 0
    }

    ol.lst-kix_list_18-1.start {
      counter-reset: lst-ctn-kix_list_18-1 0
    }

    .lst-kix_list_30-1 > li {
      counter-increment: lst-ctn-kix_list_30-1
    }

    ol.lst-kix_list_39-7.start {
      counter-reset: lst-ctn-kix_list_39-7 0
    }

    .lst-kix_list_17-7 > li:before {
      content: "" counter(lst-ctn-kix_list_17-7, decimal) ". "
    }

    ol.lst-kix_list_38-2.start {
      counter-reset: lst-ctn-kix_list_38-2 0
    }

    .lst-kix_list_33-0 > li {
      counter-increment: lst-ctn-kix_list_33-0
    }

    .lst-kix_list_17-3 > li:before {
      content: "" counter(lst-ctn-kix_list_17-3, decimal) ". "
    }

    .lst-kix_list_17-4 > li:before {
      content: "" counter(lst-ctn-kix_list_17-4, decimal) ". "
    }

    .lst-kix_list_2-8 > li:before {
      content: "" counter(lst-ctn-kix_list_2-8, decimal) ". "
    }

    .lst-kix_list_10-0 > li:before {
      content: "" counter(lst-ctn-kix_list_10-0, decimal) ". "
    }

    ol.lst-kix_list_21-7.start {
      counter-reset: lst-ctn-kix_list_21-7 0
    }

    ol.lst-kix_list_43-1.start {
      counter-reset: lst-ctn-kix_list_43-1 0
    }

    .lst-kix_list_18-3 > li:before {
      content: "" counter(lst-ctn-kix_list_18-3, decimal) ". "
    }

    .lst-kix_list_18-6 > li {
      counter-increment: lst-ctn-kix_list_18-6
    }

    ol.lst-kix_list_3-8.start {
      counter-reset: lst-ctn-kix_list_3-8 0
    }

    ol.lst-kix_list_39-8.start {
      counter-reset: lst-ctn-kix_list_39-8 0
    }

    .lst-kix_list_7-7 > li:before {
      content: "o  "
    }

    .lst-kix_list_36-5 > li:before {
      content: "" counter(lst-ctn-kix_list_36-5, lower-roman) ". "
    }

    .lst-kix_list_31-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_31-5, upper-latin) ") "
    }

    ol.lst-kix_list_20-4.start {
      counter-reset: lst-ctn-kix_list_20-4 0
    }

    .lst-kix_list_4-0 > li:before {
      content: "" counter(lst-ctn-kix_list_4-0, decimal) ". "
    }

    .lst-kix_list_36-1 > li:before {
      content: "" counter(lst-ctn-kix_list_36-1, lower-latin) ". "
    }

    .lst-kix_list_52-2 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_25-2.start {
      counter-reset: lst-ctn-kix_list_25-2 0
    }

    .lst-kix_list_15-8 > li:before {
      content: "" counter(lst-ctn-kix_list_15-8, lower-roman) ". "
    }

    ol.lst-kix_list_38-7.start {
      counter-reset: lst-ctn-kix_list_38-7 0
    }

    .lst-kix_list_15-7 > li {
      counter-increment: lst-ctn-kix_list_15-7
    }

    .lst-kix_list_4-4 > li:before {
      content: "" counter(lst-ctn-kix_list_4-4, decimal) ". "
    }

    .lst-kix_list_9-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_9-3, lower-latin) ") "
    }

    ol.lst-kix_list_13-2.start {
      counter-reset: lst-ctn-kix_list_13-2 0
    }

    ol.lst-kix_list_33-7.start {
      counter-reset: lst-ctn-kix_list_33-7 0
    }

    .lst-kix_list_29-8 > li:before {
      content: "" counter(lst-ctn-kix_list_29-8, lower-roman) ". "
    }

    ol.lst-kix_list_3-6.start {
      counter-reset: lst-ctn-kix_list_3-6 0
    }

    ol.lst-kix_list_51-5.start {
      counter-reset: lst-ctn-kix_list_51-5 0
    }

    .lst-kix_list_32-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_32-5, upper-latin) ") "
    }

    .lst-kix_list_12-4 > li:before {
      content: "" counter(lst-ctn-kix_list_12-4, decimal) " "
    }

    .lst-kix_list_5-3 > li {
      counter-increment: lst-ctn-kix_list_5-3
    }

    .lst-kix_list_33-1 > li:before {
      content: "" counter(lst-ctn-kix_list_33-1, decimal) ". "
    }

    .lst-kix_list_1-0 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) " "
    }

    ol.lst-kix_list_38-5.start {
      counter-reset: lst-ctn-kix_list_38-5 0
    }

    .lst-kix_list_11-8 > li:before {
      content: "" counter(lst-ctn-kix_list_11-8, decimal) " "
    }

    ol.lst-kix_list_2-0.start {
      counter-reset: lst-ctn-kix_list_2-0 0
    }

    .lst-kix_list_49-5 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_13-0 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) ". "
    }

    .lst-kix_list_13-4 > li:before {
      content: "" counter(lst-ctn-kix_list_13-0, decimal) "." counter(lst-ctn-kix_list_13-1, decimal) "." counter(lst-ctn-kix_list_13-2, decimal) "." counter(lst-ctn-kix_list_13-3, decimal) "." counter(lst-ctn-kix_list_13-4, decimal) ". "
    }

    ol.lst-kix_list_26-5.start {
      counter-reset: lst-ctn-kix_list_26-5 0
    }

    .lst-kix_list_55-6 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_38-6.start {
      counter-reset: lst-ctn-kix_list_38-6 0
    }

    .lst-kix_list_54-6 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) "." counter(lst-ctn-kix_list_54-2, decimal) "." counter(lst-ctn-kix_list_54-3, decimal) "." counter(lst-ctn-kix_list_54-4, decimal) "." counter(lst-ctn-kix_list_54-5, decimal) "." counter(lst-ctn-kix_list_54-6, decimal) " "
    }

    ol.lst-kix_list_2-1.start {
      counter-reset: lst-ctn-kix_list_2-1 0
    }

    ol.lst-kix_list_51-6.start {
      counter-reset: lst-ctn-kix_list_51-6 0
    }

    .lst-kix_list_4-5 > li {
      counter-increment: lst-ctn-kix_list_4-5
    }

    ol.lst-kix_list_33-8.start {
      counter-reset: lst-ctn-kix_list_33-8 0
    }

    .lst-kix_list_12-8 > li:before {
      content: "" counter(lst-ctn-kix_list_12-8, decimal) " "
    }

    .lst-kix_list_55-2 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_17-7.start {
      counter-reset: lst-ctn-kix_list_17-7 0
    }

    ol.lst-kix_list_12-6.start {
      counter-reset: lst-ctn-kix_list_12-6 0
    }

    .lst-kix_list_23-8 > li {
      counter-increment: lst-ctn-kix_list_23-8
    }

    .lst-kix_list_35-5 > li {
      counter-increment: lst-ctn-kix_list_35-5
    }

    ol.lst-kix_list_26-1.start {
      counter-reset: lst-ctn-kix_list_26-1 0
    }

    .lst-kix_list_19-2 > li:before {
      content: "" counter(lst-ctn-kix_list_19-2, decimal) ". "
    }

    ol.lst-kix_list_3-1.start {
      counter-reset: lst-ctn-kix_list_3-1 0
    }

    ol.lst-kix_list_21-0.start {
      counter-reset: lst-ctn-kix_list_21-0 0
    }

    .lst-kix_list_47-7 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_42-5.start {
      counter-reset: lst-ctn-kix_list_42-5 0
    }

    .lst-kix_list_45-0 > li {
      counter-increment: lst-ctn-kix_list_45-0
    }

    ol.lst-kix_list_51-4.start {
      counter-reset: lst-ctn-kix_list_51-4 0
    }

    .lst-kix_list_36-6 > li {
      counter-increment: lst-ctn-kix_list_36-6
    }

    .lst-kix_list_2-3 > li {
      counter-increment: lst-ctn-kix_list_2-3
    }

    .lst-kix_list_47-4 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_1-2 > li {
      counter-increment: lst-ctn-kix_list_1-2
    }

    .lst-kix_list_19-8 > li:before {
      content: "" counter(lst-ctn-kix_list_19-8, decimal) ". "
    }

    ol.lst-kix_list_20-8.start {
      counter-reset: lst-ctn-kix_list_20-8 0
    }

    .lst-kix_list_19-5 > li:before {
      content: "" counter(lst-ctn-kix_list_19-5, decimal) ". "
    }

    .lst-kix_list_47-1 > li:before {
      content: "o  "
    }

    ol.lst-kix_list_34-8.start {
      counter-reset: lst-ctn-kix_list_34-8 0
    }

    .lst-kix_list_50-2 > li {
      counter-increment: lst-ctn-kix_list_50-2
    }

    .lst-kix_list_37-7 > li:before {
      content: "" counter(lst-ctn-kix_list_37-7, lower-latin) ". "
    }

    ol.lst-kix_list_17-2.start {
      counter-reset: lst-ctn-kix_list_17-2 0
    }

    .lst-kix_list_13-2 > li {
      counter-increment: lst-ctn-kix_list_13-2
    }

    ol.lst-kix_list_21-5.start {
      counter-reset: lst-ctn-kix_list_21-5 0
    }

    .lst-kix_list_38-0 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) " "
    }

    .lst-kix_list_19-7 > li {
      counter-increment: lst-ctn-kix_list_19-7
    }

    .lst-kix_list_14-3 > li {
      counter-increment: lst-ctn-kix_list_14-3
    }

    .lst-kix_list_37-4 > li:before {
      content: "(" counter(lst-ctn-kix_list_37-4, lower-roman) ") "
    }

    .lst-kix_list_12-1 > li {
      counter-increment: lst-ctn-kix_list_12-1
    }

    .lst-kix_list_51-3 > li {
      counter-increment: lst-ctn-kix_list_51-3
    }

    .lst-kix_list_33-3 > li {
      counter-increment: lst-ctn-kix_list_33-3
    }

    ol.lst-kix_list_25-4.start {
      counter-reset: lst-ctn-kix_list_25-4 0
    }

    .lst-kix_list_38-3 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) "." counter(lst-ctn-kix_list_38-2, decimal) "." counter(lst-ctn-kix_list_38-3, decimal) " "
    }

    .lst-kix_list_38-6 > li:before {
      content: "" counter(lst-ctn-kix_list_38-0, decimal) "." counter(lst-ctn-kix_list_38-1, decimal) "." counter(lst-ctn-kix_list_38-2, decimal) "." counter(lst-ctn-kix_list_38-3, decimal) "." counter(lst-ctn-kix_list_38-4, decimal) "." counter(lst-ctn-kix_list_38-5, decimal) "." counter(lst-ctn-kix_list_38-6, decimal) " "
    }

    ol.lst-kix_list_34-3.start {
      counter-reset: lst-ctn-kix_list_34-3 0
    }

    ol.lst-kix_list_2-4.start {
      counter-reset: lst-ctn-kix_list_2-4 0
    }

    .lst-kix_list_34-4 > li {
      counter-increment: lst-ctn-kix_list_34-4
    }

    .lst-kix_list_41-8 > li {
      counter-increment: lst-ctn-kix_list_41-8
    }

    .lst-kix_list_48-0 > li:before {
      content: "\0025cf  "
    }

    ol.lst-kix_list_50-5.start {
      counter-reset: lst-ctn-kix_list_50-5 0
    }

    .lst-kix_list_2-5 > li:before {
      content: "" counter(lst-ctn-kix_list_2-5, decimal) ". "
    }

    .lst-kix_list_48-6 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_18-6 > li:before {
      content: "" counter(lst-ctn-kix_list_18-6, decimal) ". "
    }

    .lst-kix_list_39-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_39-5, upper-latin) ") "
    }

    .lst-kix_list_14-6 > li {
      counter-increment: lst-ctn-kix_list_14-6
    }

    .lst-kix_list_54-5 > li {
      counter-increment: lst-ctn-kix_list_54-5
    }

    ol.lst-kix_list_42-3.start {
      counter-reset: lst-ctn-kix_list_42-3 0
    }

    ol.lst-kix_list_39-2.start {
      counter-reset: lst-ctn-kix_list_39-2 0
    }

    .lst-kix_list_10-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_10-3, upper-latin) ") "
    }

    .lst-kix_list_15-4 > li {
      counter-increment: lst-ctn-kix_list_15-4
    }

    .lst-kix_list_2-6 > li {
      counter-increment: lst-ctn-kix_list_2-6
    }

    .lst-kix_list_36-8 > li:before {
      content: "" counter(lst-ctn-kix_list_36-8, lower-roman) ". "
    }

    .lst-kix_list_28-1 > li {
      counter-increment: lst-ctn-kix_list_28-1
    }

    .lst-kix_list_46-8 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_3-4 > li {
      counter-increment: lst-ctn-kix_list_3-4
    }

    .lst-kix_list_29-5 > li:before {
      content: "" counter(lst-ctn-kix_list_29-5, lower-latin) ") "
    }

    .lst-kix_list_20-0 > li:before {
      content: "\0025cf  "
    }

    .lst-kix_list_9-6 > li:before {
      content: "(" counter(lst-ctn-kix_list_9-6, decimal) ") "
    }

    ol.lst-kix_list_34-5.start {
      counter-reset: lst-ctn-kix_list_34-5 0
    }

    ol.lst-kix_list_51-1.start {
      counter-reset: lst-ctn-kix_list_51-1 0
    }

    .lst-kix_list_20-6 > li:before {
      content: "" counter(lst-ctn-kix_list_20-6, decimal) ". "
    }

    .lst-kix_list_23-0 > li {
      counter-increment: lst-ctn-kix_list_23-0
    }

    ol.lst-kix_list_12-1.start {
      counter-reset: lst-ctn-kix_list_12-1 0
    }

    .lst-kix_list_11-5 > li:before {
      content: "(" counter(lst-ctn-kix_list_11-5, upper-latin) ") "
    }

    ol.lst-kix_list_50-7.start {
      counter-reset: lst-ctn-kix_list_50-7 0
    }

    .lst-kix_list_54-3 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) "." counter(lst-ctn-kix_list_54-1, decimal) "." counter(lst-ctn-kix_list_54-2, decimal) "." counter(lst-ctn-kix_list_54-3, decimal) " "
    }

    ol.lst-kix_list_21-2.start {
      counter-reset: lst-ctn-kix_list_21-2 0
    }

    .lst-kix_list_20-6 > li {
      counter-increment: lst-ctn-kix_list_20-6
    }

    ol.lst-kix_list_25-7.start {
      counter-reset: lst-ctn-kix_list_25-7 0
    }

    .lst-kix_list_1-3 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) "." counter(lst-ctn-kix_list_1-2, decimal) "." counter(lst-ctn-kix_list_1-3, decimal) " "
    }

    ol.lst-kix_list_34-6.start {
      counter-reset: lst-ctn-kix_list_34-6 0
    }

    .lst-kix_list_28-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_28-3, decimal) ") "
    }

    ol.lst-kix_list_42-0.start {
      counter-reset: lst-ctn-kix_list_42-0 0
    }

    ol.lst-kix_list_2-7.start {
      counter-reset: lst-ctn-kix_list_2-7 0
    }

    .lst-kix_list_27-7 > li:before {
      content: "" counter(lst-ctn-kix_list_27-0, decimal) "." counter(lst-ctn-kix_list_27-1, decimal) "." counter(lst-ctn-kix_list_27-2, decimal) "." counter(lst-ctn-kix_list_27-3, decimal) "." counter(lst-ctn-kix_list_27-4, decimal) "." counter(lst-ctn-kix_list_27-5, decimal) "." counter(lst-ctn-kix_list_27-6, decimal) "." counter(lst-ctn-kix_list_27-7, decimal) ". "
    }

    .lst-kix_list_25-7 > li {
      counter-increment: lst-ctn-kix_list_25-7
    }

    ol.lst-kix_list_39-4.start {
      counter-reset: lst-ctn-kix_list_39-4 0
    }

    .lst-kix_list_49-2 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_35-2 > li:before {
      content: "" counter(lst-ctn-kix_list_35-0, decimal) "." counter(lst-ctn-kix_list_35-1, decimal) "." counter(lst-ctn-kix_list_35-2, decimal) " "
    }

    .lst-kix_list_3-1 > li {
      counter-increment: lst-ctn-kix_list_3-1
    }

    .lst-kix_list_30-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_30-3, lower-latin) ") "
    }

    ol.lst-kix_list_34-0.start {
      counter-reset: lst-ctn-kix_list_34-0 0
    }

    .lst-kix_list_26-8 > li {
      counter-increment: lst-ctn-kix_list_26-8
    }

    ol.lst-kix_list_39-1.start {
      counter-reset: lst-ctn-kix_list_39-1 0
    }

    .lst-kix_list_3-1 > li:before {
      content: "" counter(lst-ctn-kix_list_3-1, decimal) ". "
    }

    .lst-kix_list_14-0 > li {
      counter-increment: lst-ctn-kix_list_14-0
    }

    .lst-kix_list_44-0 > li:before {
      content: "" counter(lst-ctn-kix_list_44-0, decimal) " "
    }

    ol.lst-kix_list_17-4.start {
      counter-reset: lst-ctn-kix_list_17-4 0
    }

    .lst-kix_list_33-6 > li {
      counter-increment: lst-ctn-kix_list_33-6
    }

    ol.lst-kix_list_12-3.start {
      counter-reset: lst-ctn-kix_list_12-3 0
    }

    .lst-kix_list_44-2 > li {
      counter-increment: lst-ctn-kix_list_44-2
    }

    .lst-kix_list_21-2 > li:before {
      content: "" counter(lst-ctn-kix_list_21-2, decimal) ". "
    }

    .lst-kix_list_2-0 > li {
      counter-increment: lst-ctn-kix_list_2-0
    }

    .lst-kix_list_15-1 > li {
      counter-increment: lst-ctn-kix_list_15-1
    }

    .lst-kix_list_36-3 > li {
      counter-increment: lst-ctn-kix_list_36-3
    }

    .lst-kix_list_11-2 > li:before {
      content: "" counter(lst-ctn-kix_list_11-0, decimal) "." counter(lst-ctn-kix_list_11-1, decimal) "." counter(lst-ctn-kix_list_11-2, decimal) " "
    }

    ol.lst-kix_list_42-2.start {
      counter-reset: lst-ctn-kix_list_42-2 0
    }

    .lst-kix_list_53-1 > li:before {
      content: "o  "
    }

    .lst-kix_list_40-4 > li {
      counter-increment: lst-ctn-kix_list_40-4
    }

    ol.lst-kix_list_12-4.start {
      counter-reset: lst-ctn-kix_list_12-4 0
    }

    .lst-kix_list_35-2 > li {
      counter-increment: lst-ctn-kix_list_35-2
    }

    .lst-kix_list_25-2 > li:before {
      content: "" counter(lst-ctn-kix_list_25-2, lower-roman) ") "
    }

    .lst-kix_list_29-5 > li {
      counter-increment: lst-ctn-kix_list_29-5
    }

    .lst-kix_list_16-1 > li:before {
      content: "" counter(lst-ctn-kix_list_16-1, decimal) ". "
    }

    ol.lst-kix_list_39-0.start {
      counter-reset: lst-ctn-kix_list_39-0 0
    }

    .lst-kix_list_54-8 > li {
      counter-increment: lst-ctn-kix_list_54-8
    }

    .lst-kix_list_12-4 > li {
      counter-increment: lst-ctn-kix_list_12-4
    }

    .lst-kix_list_44-3 > li:before {
      content: "" counter(lst-ctn-kix_list_44-3, decimal) " "
    }

    .lst-kix_list_39-2 > li:before {
      content: "" counter(lst-ctn-kix_list_39-0, decimal) "." counter(lst-ctn-kix_list_39-1, decimal) "." counter(lst-ctn-kix_list_39-2, decimal) " "
    }

    .lst-kix_list_12-7 > li {
      counter-increment: lst-ctn-kix_list_12-7
    }

    .lst-kix_list_34-7 > li {
      counter-increment: lst-ctn-kix_list_34-7
    }

    .lst-kix_list_30-0 > li:before {
      content: "" counter(lst-ctn-kix_list_30-0, decimal) " "
    }

    ol.lst-kix_list_17-3.start {
      counter-reset: lst-ctn-kix_list_17-3 0
    }

    .lst-kix_list_51-0 > li {
      counter-increment: lst-ctn-kix_list_51-0
    }

    .lst-kix_list_43-4 > li:before {
      content: "" counter(lst-ctn-kix_list_43-4, decimal) " "
    }

    .lst-kix_list_7-1 > li:before {
      content: "o  "
    }

    .lst-kix_list_13-5 > li {
      counter-increment: lst-ctn-kix_list_13-5
    }

    .lst-kix_list_48-3 > li:before {
      content: "\0025aa  "
    }

    .lst-kix_list_9-6 > li {
      counter-increment: lst-ctn-kix_list_9-6
    }

    .lst-kix_list_29-2 > li {
      counter-increment: lst-ctn-kix_list_29-2
    }

    ol.lst-kix_list_42-6.start {
      counter-reset: lst-ctn-kix_list_42-6 0
    }

    .lst-kix_list_20-3 > li {
      counter-increment: lst-ctn-kix_list_20-3
    }

    .lst-kix_list_11-6 > li {
      counter-increment: lst-ctn-kix_list_11-6
    }

    .lst-kix_list_52-8 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_17-8.start {
      counter-reset: lst-ctn-kix_list_17-8 0
    }

    ol.lst-kix_list_34-1.start {
      counter-reset: lst-ctn-kix_list_34-1 0
    }

    .lst-kix_list_1-5 > li {
      counter-increment: lst-ctn-kix_list_1-5
    }

    .lst-kix_list_28-4 > li {
      counter-increment: lst-ctn-kix_list_28-4
    }

    .lst-kix_list_37-1 > li {
      counter-increment: lst-ctn-kix_list_37-1
    }

    ol.lst-kix_list_17-5.start {
      counter-reset: lst-ctn-kix_list_17-5 0
    }

    .lst-kix_list_15-2 > li:before {
      content: "(" counter(lst-ctn-kix_list_15-2, lower-latin) ") "
    }

    .lst-kix_list_10-8 > li {
      counter-increment: lst-ctn-kix_list_10-8
    }

    .lst-kix_list_10-6 > li:before {
      content: "" counter(lst-ctn-kix_list_10-6, decimal) " "
    }

    .lst-kix_list_40-7 > li {
      counter-increment: lst-ctn-kix_list_40-7
    }

    ol.lst-kix_list_12-7.start {
      counter-reset: lst-ctn-kix_list_12-7 0
    }

    .lst-kix_list_41-5 > li {
      counter-increment: lst-ctn-kix_list_41-5
    }

    .lst-kix_list_53-4 > li:before {
      content: "\0025aa  "
    }

    ol.lst-kix_list_42-8.start {
      counter-reset: lst-ctn-kix_list_42-8 0
    }

    ol.lst-kix_list_12-8.start {
      counter-reset: lst-ctn-kix_list_12-8 0
    }

    .lst-kix_list_20-3 > li:before {
      content: "" counter(lst-ctn-kix_list_20-3, decimal) ". "
    }

    .lst-kix_list_54-0 > li:before {
      content: "" counter(lst-ctn-kix_list_54-0, decimal) " "
    }

    .lst-kix_list_27-6 > li {
      counter-increment: lst-ctn-kix_list_27-6
    }

    .lst-kix_list_29-2 > li:before {
      content: "(" counter(lst-ctn-kix_list_29-2, lower-latin) ") "
    }

    .lst-kix_list_21-1 > li {
      counter-increment: lst-ctn-kix_list_21-1
    }

    ol.lst-kix_list_17-6.start {
      counter-reset: lst-ctn-kix_list_17-6 0
    }

    .lst-kix_list_28-6 > li:before {
      content: "" counter(lst-ctn-kix_list_28-6, decimal) ". "
    }

    .lst-kix_list_1-6 > li:before {
      content: "" counter(lst-ctn-kix_list_1-0, decimal) "." counter(lst-ctn-kix_list_1-1, decimal) "." counter(lst-ctn-kix_list_1-2, decimal) "." counter(lst-ctn-kix_list_1-3, decimal) "." counter(lst-ctn-kix_list_1-4, decimal) "." counter(lst-ctn-kix_list_1-5, decimal) "." counter(lst-ctn-kix_list_1-6, decimal) " "
    }

    .lst-kix_list_42-3 > li {
      counter-increment: lst-ctn-kix_list_42-3
    }

    .lst-kix_list_43-1 > li {
      counter-increment: lst-ctn-kix_list_43-1
    }

    .lst-kix_list_33-7 > li:before {
      content: "" counter(lst-ctn-kix_list_33-7, decimal) " "
    }

    ol.lst-kix_list_42-7.start {
      counter-reset: lst-ctn-kix_list_42-7 0
    }

    .lst-kix_list_40-0 > li:before {
      content: "" counter(lst-ctn-kix_list_40-0, decimal) " "
    }

    .lst-kix_list_34-3 > li:before {
      content: "(" counter(lst-ctn-kix_list_34-3, lower-latin) ") "
    }

    .lst-kix_list_2-2 > li:before {
      content: "" counter(lst-ctn-kix_list_2-2, decimal) ". "
    }

  }

</style>
