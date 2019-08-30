<template>
<div>
  <!--<sign-message id="sign-message"></sign-message>-->
  <vue-headful :title="$t('Sign up title', 'Proxeus - Sign up')"/>
  <h1 class="text-center">{{$t('Sign up')}}</h1>
  <div class="login-form container-fluid px-4 pt-2 mt-3 bg-light">
    <div class="row">
      <div class="col-6 text-center d-flex align-items-center border-right" v-if="app.wallet && metamaskLoginAvailable">
        <div class="mid align-self-center w-100">
          <h2 class="mb-3 font-weight-bold">{{$t('Wallet signature')}}</h2>
          <p class="light-text">{{$t('Use your MetaMask Wallet to sign up.')}}</p>
          <p class="text-danger" v-if="walletErrorMessage">{{ walletErrorMessage }}</p>
          <button class="btn btn-primary px-3" @click="loginWithSignature"
                  v-if="metamaskLoginAvailable">{{$t('Sign up with Metamask')}}
          </button>
        </div>
      </div>
      <div :class="{'col-6': app.wallet && metamaskLoginAvailable, 'col-12': !(app.wallet && metamaskLoginAvailable)}">
        <form v-show="!done" class="text-center" @submit.prevent="request">
          <div class="form-group mt-3 field-parent">
            <label for="inputEmail" class="sr-only">{{$t('Email address')}}</label>
            <input @input="cleanErr" type="text" id="inputEmail" ref="inputEmail" v-model.trim="email" name="email"
                   class="form-control"
                   :placeholder="$t('Email address')" required
                   autofocus>
          </div>
          <span class="text-muted"
                style="display: inline-block;">{{$t('Sign up explanation', 'Sign up by providing your email and clicking the button below.')}}</span>
          <button class="btn btn-primary px-3 mt-3" type="submit">{{$t('Sign up')}}</button>
        </form>
        <div v-show="done">
          <div
            class="my-3">{{$t('Sign up email sent explanation', 'Email sent. Please check your emails and visit the provided link to proceed.')}}
          </div>

          <a href="/" class="btn btn-primary" style="float: left;">{{$t('Home')}}</a>
          <a href="/login" class="btn btn-primary" style="float: right;">{{$t('Sign in')}}</a>
        </div>
      </div>
    </div>
  </div>
  <div class="modal fade" ref="tcModal" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{$t('Terms & Conditions')}}</h5>
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
          <button type="button" @click="acceptTC" class="btn btn-primary">{{$t('Accept')}}</button>
          <button type="button" class="btn btn-secondary" data-dismiss="modal">{{$t('Cancel')}}</button>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'RegisterRequest',
  data () {
    return {
      walletErrorMessage: '',

      account: undefined,
      email: '',
      password: '',
      metamaskLoginAvailable: false,
      loadingChallenge: false,
      challenge: null,
      done: false
    }
  },
  created () {
    if (window.web3 !== undefined) {
      this.metamaskLoginAvailable = true
    }
  },
  mounted () {
    this.$refs.inputEmail && this.$refs.inputEmail.focus()
    if (this.app.wallet) {
      this.account = this.app.wallet.getCurrentAddress()
    }
  },
  methods: {
    cleanErr () {
      $(this.$refs.inputEmail).cleanFieldErrors()
    },
    request () {
      axios.post('/api/register', { email: this.email }).then(res => {
        this.cleanErr()
        this.done = true
      }, (err) => {
        this.cleanErr()
        console.log(err)
        this.app.handleError(err)
        if (err.response && err.response.status === 422) {
          $(this.$refs.inputEmail).showFieldErrors({ errors: err.response.data })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Warning'),
            text: this.$t('There was an unexpected error. Please try again or if the error persists contact the platform operator.'),
            type: 'warning'
          })
        }
        this.$nextTick(() => {
          this.$refs.inputEmail.focus()
        })
      })
    },
    acceptTC () {
      localStorage.setItem('acc_' + this.account, 'yes')
      this.app.acknowledgeFirstLogin()
      $(this.$refs.tcModal).modal('hide')
      this.metamaskLogin()
    },
    checkTermsAndConditions () {
      let rememberAccept = localStorage.getItem('acc_' + this.account)
      if (rememberAccept && rememberAccept === 'yes') {
        return true
      }
      $(this.$refs.tcModal).modal('show')
    },
    loginWithSignature () {
      if (!this.challenge) {
        if (window.web3 !== undefined) {
          this.metamaskLoginAvailable = true
          axios.get('/api/challenge').then((response) => {
            this.challenge = response.data
            this.metamaskLogin()
          }, (err) => {
            this.app.handleError(err)
          })
        }
      } else {
        this.metamaskLogin()
      }
    },
    async metamaskLogin () {
      if (window.ethereum) {
        try {
          await window.ethereum.enable()
          await this.app.wallet.wallet.setupDefaultAccount()
        } catch (e) {
          console.log(e)
          this.walletErrorMessage = this.$t('Please grant access to MetaMask.')
          return
        }
      } else {
        this.walletErrorMessage = this.$t('Please grant access to MetaMask.')
        return
      }
      this.account = this.app.wallet.getCurrentAddress()
      if (this.account === undefined) {
        this.walletErrorMessage = this.$t('Please sign in to MetaMask.')
        return
      }
      if (this.checkTermsAndConditions()) {
        this.app.wallet.signMessage(this.challenge, this.account).then((signature) => {
          axios.post('/api/login', { signature }).then((res) => {
            this.challenge = ''
            if (res.status >= 200 && res.status <= 299) {
              window.location = res.data.location || '/admin/workflow'
            } else {
              this.walletErrorMessage = this.$t('Could not verify signature.')
            }
          }, (err) => {
            this.challenge = ''
            this.app.handleError(err)
            this.walletErrorMessage = this.$t('Could not verify signature.')
          })
        }).catch(() => {
          this.walletErrorMessage = this.$t('Could not Sign Message.')
        })
      }
    }
  }
}
</script>

<style lang="scss">
  @import "../assets/styles/variables.scss";

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
</style>
