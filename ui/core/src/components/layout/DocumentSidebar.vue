<template>
<nav class="col-md-3 col-lg-2 d-none d-md-block app-sidebar px-0 pt-0" :class="{toggled:toggled}">
  <a class="navbar-brand mb-1 w-100 p-3 mx-0" href="#">
    <img src="../../assets/images/logo.png" width="30" height="30" alt="" class="d-inline-block align-top">
    <span class="ml-1 brand-name">Brand name</span>
  </a>
  <nav class="collapse show sidebar-sticky">
    <nav class="nav main-nav flex-column">
      <li class="nav-item">
        <router-link :to="{name:'Workflows'}" class="nav-link"><span
          class="material-icons">share</span><span class="nav-link-title">Workflows</span>
        </router-link>
      </li>
      <li class="nav-item">
        <router-link :to="{name:'Templates'}" class="nav-link"><span class="material-icons">web</span><span
          class="nav-link-title">Templates</span>
        </router-link>
      </li>
      <li class="nav-item">
        <router-link :to="{name:'Forms'}" class="nav-link"><span
          class="material-icons">view_quilt</span><span class="nav-link-title">Forms</span>
        </router-link>
      </li>
      <li class="nav-item">
        <router-link :to="{name:'Users'}" class="nav-link"><span class="material-icons">people</span><span
          class="nav-link-title">Users</span>
        </router-link>
      </li>
      <li class="nav-item">
        <router-link :to="{name:'I18n'}" class="nav-link"><span class="material-icons">language</span><span
          class="nav-link-title">Localization</span>
        </router-link>
      </li>
    </nav>
  </nav>
</nav>
</template>

<script>
export default {
  name: 'sidebar',
  props: ['user', 'toggled'],
  mounted () {
    $('[data-toggle="tooltip"]').tooltip()
  }
}
</script>

<style lang="scss">
  .app-sidebar {
    background: #062a85;
    position: -webkit-sticky;
    position: sticky;
    top: 0;
    height: calc(100vh);
    padding-top: .5rem;
    overflow-x: hidden;
    overflow-y: scroll;
    max-width: 230px;
    -webkit-transition: all .2s ease;
    transition: all .2s ease;

    .nav-link {
      vertical-align: middle;
      color: #cedfff;

      &.dropdown-toggle::after {
        margin-left: 1rem;
        vertical-align: middle;
      }

      &.router-link-exact-active {
        background: #1a3d96;
      }

      &:hover {
        background: #254394;
      }

      span {
        vertical-align: middle;
      }
    }

    .nav-link .material-icons {
      font-size: 1.6rem;
      margin-right: 30px;
      vertical-align: middle;
      color: white;

    }

    .nav-item {
      width: 100%;
    }

    &.toggled {
      max-width: 80px;

      .nav-item {
        width: 100%;
      }

      &:hover {
        .brand-name {
          display: inline-block;
        }

        .nav-link-title {
          display: inline-block;
        }

        .nav-link.dropdown-toggle::after {
          display: inline-block;
        }
      }

      .brand-name {
        display: none;
      }

      .nav-link-title {
        display: none;
      }

      .nav-link.dropdown-toggle::after {
        display: none;
      }

    }

    .navbar-brand {
      color: #e5efff;
    }
    .secondary-nav {
      position: fixed;
      bottom: 0;
      width: 100%;
      max-width: 230px;

      .nav-item {
        width: 100%;
        max-width: 230px;
      }
    }
    .main-nav {
      height: 100%;
      margin-bottom: 150px !important;
    }
  }
</style>
