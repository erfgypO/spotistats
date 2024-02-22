// Composables
import { createRouter, createWebHistory } from 'vue-router'
import {useAppStore} from "@/store/app";
import {useUserStore} from "@/store/user";

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import('@/views/Home.vue'),
      },
      {
        path: 'connect',
        name: 'Connect',
        component: () => import('@/views/ConnectSpotify.vue'),
      },
      {
        path: 'account',
        name: 'Account',
        component: () => import('@/views/Account.vue'),
      }
    ],
  },
  {
    path: '/anon',
    component: () => import('@/layouts/default/Anon.vue'),
    children: [
      {
        path: 'auth',
        name: 'Auth',
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import('@/views/Auth.vue'),
      },
    ]
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

router.beforeEach((to, from, next) => {
  const appStore = useAppStore();
  if(!to.path.startsWith('/anon') && (appStore.expiresAt * 1000 <= Date.now() || localStorage.getItem('token') === null)){
    next('/anon/auth');
    appStore.expiresAt = 0;
    appStore.accessToken = "";

    return;
  }

  const userStore = useUserStore();
  if(!to.path.startsWith('/anon') && to.name !== 'Connect' && !userStore.connectedToSpotify) {
    next('/connect');
    return;
  }
  next();
});

export default router
