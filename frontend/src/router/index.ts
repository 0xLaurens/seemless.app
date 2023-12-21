import {createRouter, createWebHistory} from 'vue-router'
import UsernameView from '../views/UsernameView.vue'
import ErrorView from '@/views/ErrorView.vue'
import ExchangeView from "@/views/ExchangeView.vue";
import QrView from "@/views/QrView.vue";
import {useWebsocketStore} from "@/stores/websocket";
import SettingsView from "@/views/SettingsView.vue";
import AboutView from "@/views/AboutView.vue";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: ExchangeView
        },
        {
            path: '/share',
            name: 'share',
            component: QrView,
        },
        {
            path: '/settings',
            name: 'settings',
            component: SettingsView,
        },
        {
            path: '/about',
            name: 'about',
            component: AboutView,
        },
        {
            path: '/nick',
            name: 'username',
            component: UsernameView
        },
        {
            path: '/404',
            name: '404 not found',
            component: ErrorView
        },
        {
            path: '/:unknown*',
            redirect: '/404'
        }
    ]
})

router.beforeEach((to, _, next) => {
    if (to.name === 'home' || to.name === '404 not found') {
        next()
    } else {
        const ws = useWebsocketStore()

        if (ws.isClosed()) {
            return next({name: 'home'})
        }
        return next()
    }
})

// router.beforeEach((to, from, next) => {
//     const store = useUserStore()
//     const emptyUsername = store.getUsername() == ''
//     if ((to.name == 'room' && emptyUsername) || (to.name == 'network' && emptyUsername)) {
//         next({name: 'username'})
//     } else {
//         next()
//     }
// })

export default router
