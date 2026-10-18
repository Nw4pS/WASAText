import { createRouter, createWebHashHistory } from 'vue-router'
// importiamo le due pagine principali
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
      // meta è lo spazio dove attaccare delle etichette alle rotte
      // qui diciamo che questa pagina ha bisogno dell'autenticazione.
      meta: { richiedeLogin: true }
    }
  ]
})

// "beforeEach" scatta prima di ogni navigazione.
// "to" è la pagina dove sto andando, "from" da dove vengo, "next" è il semaforo (vai avanti).
router.beforeEach((to, from, next) => {
  // cerco il token nel localS1torage del browser
  const isLoggato = localStorage.getItem('token')

  // caso 1, l'utente vuole andare in una pagina protetta ma non è loggato
  if (to.meta.richiedeLogin && !isLoggato) {
    next('/login')
  } 
  // caso 2, l'utente è già loggato e prova a tornare alla schermata di login
  else if (to.path === '/login' && isLoggato) {
    next('/')
  } 
  // caso 3, tutto ok, lo lascio passare normalmente
  else {
    next()
  }
})

export default router