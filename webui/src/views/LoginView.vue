<template>
  <!-- contenitore principale per centrare tutto nella pagina -->
  <div class="login-wrapper">
    <div class="login-card">
      <h2>benvenuto su wasatext</h2>
      
      <!-- form di login. .prevent evita che la pagina si ricarichi quando premo invio -->
      <form @submit.prevent="faiLogin">
        <div class="input-group">
          <label>nome utente</label>
          <!-- v-model collega quello che scrivo alla variabile 'username' -->
          <!-- metto i limiti di 3 e 16 caratteri come chiede l'api.yaml -->
          <input 
            type="text" 
            v-model="username" 
            minlength="3" 
            maxlength="16" 
            required 
            placeholder="es. mario"
          >
        </div>
        
        <!-- bottone che si disabilita se sta caricando -->
        <button type="submit" :disabled="caricamento">
          {{ caricamento ? 'sto entrando...' : 'entra' }}
        </button>
      </form>

      <!-- se c'è un errore, mostro questo paragrafo -->
      <p v-if="errore" class="errore-testo">{{ errore }}</p>
    </div>
  </div>
</template>

<script>
// importo axios che abbiamo configurato prima per fare le chiamate al backend
import api from '../services/axios.js';

export default {
  data() {
    return {
      // qui salvo il nome che l'utente scrive nell'input
      username: '',
      // mi serve per far capire all'utente che sto aspettando il server
      caricamento: false,
      // se qualcosa va male, scrivo l'errore qui
      errore: ''
    };
  },
  methods: {
    // funzione asincrona perché deve aspettare la risposta del server
    async faiLogin() {
      // inizio a caricare e pulisco vecchi errori
      this.caricamento = true;
      this.errore = '';
      
      try {
        // faccio la chiamata post a /session passando il nome
        const risposta = await api.post('/session', {
          name: this.username
        });
        
        // il backend mi restituisce un identifier
        const identificatore = risposta.data.identifier;
        
        // salvo l'identificatore come token nel browser per usarlo nelle altre chiamate
        localStorage.setItem('token', identificatore);
        // salvo anche il nome così so chi sono
        localStorage.setItem('username', this.username);
        
        // tutto è andato bene, mi sposto sulla pagina principale della chat
        this.$router.push('/');
        
      } catch (e) {
        // se il server dà errore (es. non è acceso), lo salvo nella variabile
        this.errore = 'impossibile fare il login. controlla se il backend è acceso.';
      } finally {
        // ho finito, sia se è andata bene sia se è andata male
        this.caricamento = false;
      }
    }
  }
};
</script>

<style scoped>
/* un po' di stile per rendere il box di login carino e centrato */
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}

.login-card {
  background: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  width: 300px;
  text-align: center;
}

.input-group {
  margin-bottom: 20px;
  text-align: left;
}

.input-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.input-group input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background-color: #ccc;
}

.errore-testo {
  color: red;
  margin-top: 15px;
  font-size: 14px;
}
</style>