import {createRouter, createWebHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UsernameView from '../views/UsernameView.vue'
import NetworkView from '@/views/NetworkView.vue'
import RoomView from '@/views/RoomView.vue'
import {useUserStore} from '@/stores/user'
import ErrorView from '@/views/ErrorView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView
        },
        {
            path: '/nick',
            name: 'username',
            component: UsernameView
        },
        {
            path: '/network',
            name: 'network',
            component: NetworkView
        },
        {
            path: '/room/:id',
            name: 'room',
            component: RoomView
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
router.beforeEach((to, from, next) => {
    const store = useUserStore()
    const emptyUsername = store.getUsername() == ''
    if ((to.name == 'room' && emptyUsername) || (to.name == 'network' && emptyUsername)) {
        next({name: 'username'})
    } else {
        next()
    }
})

export default router
