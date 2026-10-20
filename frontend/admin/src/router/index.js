// Composables
import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/store/auth";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import(/* webpackChunkName: "login" */ "@/views/Login.vue"),
  },
  {
    path: "/",
    component: () => import("@/layouts/default/Default.vue"),
    meta: { transition: "slide-right" },
    children: [
      {
        path: "",
        name: "Users",
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () =>
          import(/* webpackChunkName: "users" */ "@/views/Users.vue"),
      },
    ],
  },
  {
    path: "/user-profile/:id",
    component: () => import("@/layouts/default/Default.vue"),
    meta: { transition: "slide-right" },
    children: [
      {
        path: "",
        name: "UserProfile",
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () =>
          import(
            /* webpackChunkName: "userProfile" */ "@/views/UserProfile.vue"
          ),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

// Only the very first navigation needs to hit the server - after that we
// trust the in-memory isAuthenticated flag that login()/logout() keep
// up to date, so we don't re-verify on every single route change.
let authChecked = false;

router.beforeEach(async (to) => {
  const authStore = useAuthStore();

  if (!authChecked) {
    authChecked = true;
    await authStore.checkAuth();
  }

  if (to.name !== "Login" && !authStore.isAuthenticated) {
    return { name: "Login", query: { redirect: to.fullPath } };
  }

  if (to.name === "Login" && authStore.isAuthenticated) {
    return { name: "Users" };
  }
});

export default router;
