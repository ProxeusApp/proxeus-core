<template>
<div class="view-support">
  <vue-headful :title="$t('Help / Support title', 'Proxeus - Help / Support')"/>
  <top-nav :title="$t('Help / Support')">
  </top-nav>
  <div class="container-fluid" style="padding-top: 1.5rem;">
    <div class="row mb-3">
      <div class="col-sm-12">
        <section>
          <a target="_blank" href="https://proxeus.com/en/beta/library/"
             class="btn btn-primary">{{$t('Documentation Library')}}
          </a>
          <a target="_blank"
             href="https://docs.google.com/document/d/1nUqB31uWRVokDOeKgzhv49rTbkLQAS8nmv8bCumDMEM/preview"
             class="btn btn-secondary ml-2">{{$t('Handbook')}}
          </a>
        </section>
        <hr>
        <p
          class="light-text">{{$t('Please use the form below to provide us with your feedback or support request.')}}</p>
        <iframe title="Feedback Form" class="freshwidget-embedded-form" id="freshwidget-embedded-form"
                src="https://blockfactory.freshdesk.com/widgets/feedback_widget/new?&widgetType=embedded&formTitle=Contact+Form&submitTitle=Submit&submitThanks=Thank+you!+We+well+received+your+message+and+will+get+back+to+you+as+soon+as+possible.&screenshot=No"
                scrolling="no" height="500px" width="100%" frameborder="0">
        </iframe>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import TopNav from '../components/layout/TopNav'

export default {
  name: 'support',
  components: {
    TopNav
  }
}
</script>

<style scoped>
  @import url("https://s3.amazonaws.com/assets.freshdesk.com/widget/freshwidget.css");

  .freshwidget-button a:hover, .freshwidget-theme {
    text-decoration: none
  }

  .widget-ovelay {
    filter: alpha(opacity=50);
    opacity: .5;
    background-color: #ffffff
  }

  .freshwidget-container, .widget-ovelay {
    position: fixed;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    z-index: 90002
  }

  .fd-btn-left a {
    left: -2px
  }

  .fd-btn-left a:hover {
    left: 0
  }

  .fd-btn-right {
    right: -2px
  }

  .fd-btn-right:hover {
    right: 0
  }

  .ie8.fd-btn-right, .ie9.fd-btn-right {
    right: 35px
  }

  .ie8.fd-btn-right:hover, .ie9.fd-btn-right:hover {
    right: 33px
  }

  .fd-btn-top {
    top: -2px
  }

  .fd-btn-top:hover {
    top: 0
  }

  .fd-btn-bottom {
    bottom: 0
  }

  .fd-btn-bottom a {
    bottom: -2px
  }

  .fd-btn-bottom a:hover {
    bottom: 0
  }

  .freshwidget-theme {
    text-shadow: rgba(0, 0, 0, .2), -1px, -1px, 1px;
    padding: 8px 16px;
    background-color: #015453;
    color: #ffffff;
    font-size: 14px;
    font-weight: 700;
    outline: 0
  }

  .freshwidget-theme:hover {
    color: #f3f3f3
  }

  .freshwidget-customimage {
    overflow: visible
  }

  .fd-btn-right .freshwidget-customimage {
    right: 0
  }

  .freshwidget-button {
    position: fixed;
    display: block !important;
    width: 0;
    height: 0;
    z-index: 90000;
    margin: 0;
    outline: 0
  }

  .freshwidget-button a {
    position: fixed;
    white-space: nowrap;
    outline: 0
  }

  .proxy-link {
    position: absolute;
    z-index: 10;
    outline: 0
  }

  .fd-btn-left a.freshwidget-theme, .fd-btn-top a.freshwidget-theme {
    -moz-border-radius-bottomleft: 6px;
    -webkit-border-bottom-left-radius: 6px;
    border-bottom-left-radius: 6px;
    -moz-border-radius-bottomright: 6px;
    -webkit-border-bottom-right-radius: 6px;
    border-bottom-right-radius: 6px;
    -webkit-box-shadow: 0 0 10px rgba(0, 0, 0, .1), inset 0 -1px 0 rgba(0, 0, 0, .08), inset 0 -2px 0 rgba(255, 255, 255, .7), inset 0 9px 29px rgba(0, 0, 0, .1);
    -moz-box-shadow: 0 0 10px rgba(0, 0, 0, .1), inset 0 -1px 0 rgba(0, 0, 0, .08), inset 0 -2px 0 rgba(255, 255, 255, .7), inset 0 9px 29px rgba(0, 0, 0, .1);
    box-shadow: 0 0 10px rgba(0, 0, 0, .1), inset 0 -1px 0 rgba(0, 0, 0, .08), inset 0 -2px 0 rgba(255, 255, 255, .7), inset 0 9px 29px rgba(0, 0, 0, .1)
  }

  .fd-btn-left a.freshwidget-theme {
    transform: rotate(-90deg);
    transform-origin: top left;
    -moz-transform: rotate(-90deg);
    -moz-transform-origin: top left;
    -webkit-transform: rotate(-90deg);
    -webkit-transform-origin: top left;
    -o-transform: rotate(-90deg);
    -o-transform-origin: top left;
    filter: progid:DXImageTransform.Microsoft.BasicImage(rotation=3)
  }

  .fd-btn-bottom a.freshwidget-theme, .fd-btn-right a.freshwidget-theme {
    -moz-border-radius-topleft: 6px;
    -webkit-border-top-left-radius: 6px;
    border-top-left-radius: 6px;
    -moz-border-radius-topright: 6px;
    -webkit-border-top-right-radius: 6px;
    border-top-right-radius: 6px;
    -webkit-box-shadow: 0 0 10px rgba(0, 0, 0, .1), inset 0 1px 0 rgba(0, 0, 0, .08), inset 0 2px 0 rgba(255, 255, 255, .7), inset 0 -9px 29px rgba(0, 0, 0, .1);
    -moz-box-shadow: 0 0 10px rgba(0, 0, 0, .1), inset 0 1px 0 rgba(0, 0, 0, .08), inset 0 2px 0 rgba(255, 255, 255, .7), inset 0 -9px 29px rgba(0, 0, 0, .1);
    box-shadow: 0 0 10px rgba(0, 0, 0, .1), inset 0 1px 0 rgba(0, 0, 0, .08), inset 0 2px 0 rgba(255, 255, 255, .7), inset 0 -9px 29px rgba(0, 0, 0, .1)
  }

  .fd-btn-right a.freshwidget-theme {
    transform: rotate(-90deg);
    transform-origin: bottom left;
    -moz-transform: rotate(-90deg);
    -moz-transform-origin: bottom left;
    -webkit-transform: rotate(-90deg);
    -webkit-transform-origin: bottom left;
    -o-transform: rotate(-90deg);
    -o-transform-origin: bottom left;
    filter: progid:DXImageTransform.Microsoft.BasicImage(rotation=3)
  }

  .fd-btn-left a.ie10 {
    -ms-transform: widget-rotate(-90deg);
    -ms-transform-origin: top left
  }

  .fd-btn-right a.ie10 {
    -ms-transform: widget-rotate(-90deg);
    -ms-transform-origin: bottom left
  }

  .freshwidget-dialog {
    position: absolute;
    padding: 8px 0;
    left: 50%;
    margin-left: -260px;
    top: 20px;
    z-index: 900010
  }

  @-moz-keyframes freshwidget-dialog-box-animation {
    from {
      margin-left: -260px
    }
    to {
      margin-left: -420px
    }
  }

  @-webkit-keyframes freshwidget-dialog-box-animation {
    from {
      margin-left: -260px
    }
    to {
      margin-left: -420px
    }
  }

  @-o-keyframes freshwidget-dialog-box-animation {
    from {
      margin-left: -260px
    }
    to {
      margin-left: -420px
    }
  }

  @keyframes freshwidget-dialog-box-animation {
    from {
      margin-left: -260px
    }
    to {
      margin-left: -420px
    }
  }

  @-moz-keyframes iframeResize-animation {
    from {
      width: 500px
    }
    to {
      width: 850px
    }
  }

  @-webkit-keyframes iframeResize-animation {
    from {
      width: 500px
    }
    to {
      width: 850px
    }
  }

  @-o-keyframes iframeResize-animation {
    from {
      width: 500px
    }
    to {
      width: 850px
    }
  }

  @keyframes iframeResize-animation {
    from {
      width: 500px
    }
    to {
      width: 850px
    }
  }

  @-moz-keyframes freshwidget-dialog-box-animation-close {
    from {
      margin-left: -420px
    }
    to {
      margin-left: -260px
    }
  }

  @-webkit-keyframes freshwidget-dialog-box-animation-close {
    from {
      margin-left: -420px
    }
    to {
      margin-left: -260px
    }
  }

  @-o-keyframes freshwidget-dialog-box-animation-close {
    from {
      margin-left: -420px
    }
    to {
      margin-left: -260px
    }
  }

  @keyframes freshwidget-dialog-box-animation-close {
    from {
      margin-left: -420px
    }
    to {
      margin-left: -260px
    }
  }

  @-moz-keyframes iframe-animation-close {
    from {
      width: 850px
    }
    to {
      width: 526px
    }
  }

  @-webkit-keyframes iframe-animation-close {
    from {
      width: 850px
    }
    to {
      width: 526px
    }
  }

  @-o-keyframes iframe-animation-close {
    from {
      width: 850px
    }
    to {
      width: 526px
    }
  }

  @keyframes iframe-animation-close {
    from {
      width: 850px
    }
    to {
      width: 526px
    }
  }

  .freshwidget-dialog-box-close {
    animation: freshwidget-dialog-box-animation-close .15s ease-in-out forwards
  }

  .freshwidget-dialog-box {
    animation: freshwidget-dialog-box-animation .15s ease-in-out forwards
  }

  .freshwidget-dialog iframe {
    width: 526px;
    height: 600px
  }

  .freshwidget-dialog iframe.iframeClose {
    animation: iframe-animation-close .15s ease-in-out forwards
  }

  .freshwidget-dialog iframe.iframeResize {
    animation: iframeResize-animation .15s ease-in-out forwards
  }

  .freshwidget-dialog .widget-close {
    display: block;
    position: absolute;
    left: -12px;
    top: 16px;
    width: 32px;
    height: 32px;
    background: #6d7574;
    cursor: pointer;
    color: #ffffff;
    border: 1px solid #ffffff;
    border-radius: 3px 0 0 3px;
    border-right-color: #6d7574;
    font-size: 26px;
    font-style: normal;
    font-weight: 600;
    line-height: 20px;
    font-family: HelveticaNeue-Light, "Helvetica Neue Light", "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif
  }

  .widget-close:before {
    content: "\00d7";
    position: absolute;
    left: 9px;
    top: 4px
  }

  .mobile-widget-arrow, .mobile-widget-close {
    display: none
  }

  @media (max-width: 768px) {
    .responsive .mobile-widget-arrow-block, .responsive .mobile-widget-close {
      display: block
    }

    .responsive .freshwidget-container, .responsive .widget-ovelay {
      background: #ffffff
    }

    .responsive .mobile-widget-arrow, .responsive .mobile-widget-close {
      position: fixed;
      top: 2%;
      font-size: 32px;
      color: #000000;
      font-weight: 700;
      margin: 0;
      padding: 0;
      left: 15px;
      z-index: 10;
      cursor: pointer;
      line-height: 15px;
      font-family: HelveticaNeue-Light, "Helvetica Neue Light", "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif
    }

    .responsive .freshwidget-dialog .widget-close, .responsive .mobile-widget-close-none {
      display: none
    }

    .responsive .mobile-widget-arrow:before {
      content: "\003C"
    }

    .responsive .mobile-widget-close:before {
      content: "\00d7"
    }

    .responsive .freshwidget-dialog {
      left: 0;
      margin-left: 0;
      top: 0
    }

    .responsive .freshwidget-dialog iframe {
      position: fixed;
      top: 0;
      height: 100% !important;
      width: 100%
    }
  }
</style>
