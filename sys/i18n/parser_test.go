package i18n

import (
	"log"
	"testing"
)

var str = `$t('verification.blockchain.creator', 'Issued by')
$t('verification.blockchain.notFound',
              'No entry was found for this file. This could mean it was never registered or it was manipulated. If you believe that this finding is an error, please contact the issuer.')
$t('verification.blockchain.invalid',
              'This file has been recognised as authentic but has since been declared invalid. It may have expired or been recalled by the issuer.')
$t('verification.block)chain.error', 'File could not been verified due to technical problems. Please try again.')


$t('verification.block)chain.error')

            </div>
            </div>
            <h2 class="mt-3">{{ $t("backend.internationalization.title.languages") }}</h2>
            <div class="row">
      <div v-if="!loading">
        <div class="break-word mx-auto mt-2 mb-0">
          <h4 v-show="valid" class="text-success">{{ $t('verification.blockchain.hint.valid', 'The file {filename} is valid.', {filename: file.name}) }}</h4>
          <h4 v-show="invalid" class="text-danger">{{ $t('verification.blockchain.hint.invalid', 'The file {filename} has been revoked.', {filename: file.name}) }}</h4>
          <h4 v-show="notFound" class="text-danger">{{ $t('verification.blockchain.hint.notFound', 'The file {filename} is invalid.', {filename: file.name}) }}</h4>
        </div>
        <div class="text-center mt-2">

      <div class="animscroll">
        <spinner v-show="initialized === false" background="transparent" color="#eee"></spinner>
        <div class="light-text mx-auto mt-5" v-show="initialized === false">{{ $t('verification.blockchain.connecting', 'Connecting to network ...') }}</div>
        <div class="dropbox bg-white pl-0" v-if="initialized"
             :class="{hasDocuments}">
          <!-- todo: on file drag expand dropzone -->
          <file-drop-box @dropped="drop"></file-drop-box>
          <p class="w-100 text-center" v-if="!hasDocuments">
            <span class="text-hint py-2 px-4 d-inline-block">
             {{ $t('verification.blockchain.verify', 'You can verify the authenticity of your document here. The hash of your document will be compared to the hashes, which were registered on the blockchain upon document creation.') }}
            </span>
          </p>
        </div>
        <div class="stats col-sm-12 p-0" v-show="hasDocuments">
            <verification-file-entry :singleFile="files.length === 1" v-for="file in files" :key="file.name" :file="file" :wallet="wallet"
                                     :thumbnail="thumbnail"></verification-file-entry>
        </div>
      </div>
            <div v-if="meta && meta.langList && meta.langList.length>0">
            <h2 style="margin-top:20px;">{{ $t('backend.internationalization.title.translation') }}</h2>
            <div style="position:relative;">`

func TestI18n_Close(t *testing.T) {
	p := NewUIParser()
	p.Parse([]byte(str))
	trans := p.Translations()
	for k, v := range trans {
		log.Println(k, "::::", v)
	}
}
