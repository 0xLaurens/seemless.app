import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UsernameView from '../views/UsernameView.vue'
import NetworkView from '@/views/NetworkView.vue'
import RoomView from '@/views/RoomView.vue'

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
    }
  ]
})

export default router
