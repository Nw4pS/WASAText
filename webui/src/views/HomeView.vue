<template>
  <div class="container-fluid vh-100 d-flex flex-column p-0">
    <!-- Header superiore -->
    <header class="bg-primary text-white p-3 d-flex justify-content-between align-items-center">
      <h5 class="m-0">WASAText</h5>
      <div>
        <span class="small me-3">Utente: {{ mioUsername }}</span>
        <button @click="mostraModaleGruppo = true" class="btn btn-sm btn-light me-2">Nuovo Gruppo</button>
        <button @click="logout" class="btn btn-sm btn-outline-light">Esci</button>
      </div>
    </header>

    <div class="row flex-grow-1 g-0 overflow-hidden">
      <!-- Sidebar sinistra: 4 colonne -->
      <aside class="col-4 border-end bg-light d-flex flex-column h-100">
        
        <!-- Barra di ricerca -->
        <div class="p-3 border-bottom bg-white">
          <div class="input-group">
            <input 
              type="text" 
              class="form-control" 
              placeholder="Cerca utente per iniziare chat..." 
              v-model="testoRicerca"
              @focus="attivaModalitaRicerca"
              @input="cercaUtenti"
            />
            <!-- Tasto per uscire dalla modalità ricerca e tornare alle proprie chat -->
            <button 
              v-if="modalitaRicerca" 
              @click="annullaRicerca" 
              class="btn btn-outline-secondary" 
              type="button"
            >
              &times;
            </button>
          </div>
        </div>

        <div class="list-group list-group-flush overflow-auto flex-grow-1">
          
          <!-- SEZIONE RICERCA / ELENCO TUTTI GLI UTENTI -->
          <div v-if="modalitaRicerca">
            <div class="p-2 bg-light text-muted small fw-bold d-flex justify-content-between align-items-center">
              <span>{{ testoRicerca.trim() === '' ? 'Tutti gli utenti registrati:' : 'Risultati ricerca:' }}</span>
              <button @click="annullaRicerca" class="btn btn-link btn-sm p-0 text-decoration-none">Torna alle chat</button>
            </div>

            <div v-if="caricamentoUtenti" class="p-4 text-center text-muted">
              Caricamento contatti...
            </div>

            <div v-else-if="utentiTrovati.length === 0" class="p-3 text-center text-muted">
              Nessun utente trovato.
            </div>

            <!-- Lista utenti cliccabili per avviare chat diretta (POST /conversations) -->
            <button
              v-else
              v-for="utente in utentiTrovati"
              :key="utente.username"
              @click="avviaChatDiretta(utente.username)"
              class="list-group-item list-group-item-action d-flex align-items-center p-3"
            >
              <img
                v-if="utente.proPic"
                :src="utente.proPic"
                class="rounded-circle me-3"
                style="width: 40px; height: 40px; object-fit: cover;"
              />
              <div 
                v-else 
                class="rounded-circle bg-secondary text-white d-flex justify-content-center align-items-center me-3" 
                style="width: 40px; height: 40px; font-weight: bold;"
              >
                {{ utente.username.charAt(0).toUpperCase() }}
              </div>
              <div class="flex-grow-1 text-start">
                <strong>{{ utente.username }}</strong>
                <div class="small text-muted">Clicca per avviare la chat</div>
              </div>
            </button>
          </div>

          <!-- SEZIONE CHAT ESISTENTI (VISIBILE QUANDO NON SI CERCA) -->
          <div v-else>
            <div class="p-2 bg-light text-muted small fw-bold">Le mie chat:</div>
            
            <div v-if="caricamentoChat" class="p-4 text-center text-muted">
              Caricamento conversazioni...
            </div>

            <ConversationItem 
              v-else
              v-for="conv in conversazioni" 
              :key="conv.id" 
              :conv="conv" 
              @seleziona="apriChat(conv.id)" 
            />

            <div v-if="conversazioni.length === 0 && !caricamentoChat" class="p-4 text-center text-muted">
              Nessuna chat attiva. Clicca sulla barra di ricerca in alto per trovare un utente!
            </div>
          </div>

        </div>
      </aside>

      <!-- Area principale destra: 8 colonne -->
      <main class="col-8 d-flex h-100 bg-white justify-content-center align-items-center">
        <div class="text-center text-muted">
          <h4>Seleziona una chat</h4>
          <p>Scegli una conversazione a sinistra o cerca un utente per scrivergli.</p>
        </div>
      </main>
    </div>

    <!-- Modale creazione nuovo gruppo -->
    <NewGroupModal 
      v-if="mostraModaleGruppo" 
      :mioUsername="mioUsername" 
      @chiudi="mostraModaleGruppo = false" 
      @fatto="gestisciGruppoCreato" 
    />
  </div>
</template>

<script>
import api from '../services/axios.js';
import ConversationItem from '../components/ConversationItem.vue';
import NewGroupModal from '../components/NewGroupModal.vue';

export default {
  components: {
    ConversationItem,
    NewGroupModal
  },
  data() {
    return {
      mioUsername: localStorage.getItem('username') || '',
      conversazioni: [],
      caricamentoChat: false,

      // Gestione Ricerca / Visualizzazione utenti
      modalitaRicerca: false,
      testoRicerca: '',
      utentiTrovati: [],
      caricamentoUtenti: false,

      mostraModaleGruppo: false
    };
  },
  methods: {
    // 1. Scarica le conversazioni dell'utente (GET /conversations)
    async scaricaConversazioni() {
      this.caricamentoChat = true;
      try {
        const risposta = await api.get('/conversations');
        let chatRicevute = risposta.data || [];
        
        // Ordine cronologico inverso come da specifiche PDF
        chatRicevute.sort((a, b) => {
          if (!a.lastMessageDate) return 1;
          if (!b.lastMessageDate) return -1;
          return new Date(b.lastMessageDate) - new Date(a.lastMessageDate);
        });
        
        this.conversazioni = chatRicevute;
      } catch (e) {
        console.error("Errore durante il download delle conversazioni", e);
      } finally {
        this.caricamentoChat = false;
      }
    },

    // 2. Attiva la modalità ricerca al click sulla barra e scarica tutti gli utenti
    async attivaModalitaRicerca() {
      this.modalitaRicerca = true;
      // Se non c'è testo digitato, scarica tutti gli utenti registrati
      if (this.testoRicerca.trim() === '') {
        await this.eseguiChiamataUtenti('');
      }
    },

    // 3. Esegue la ricerca digitando testo
    async cercaUtenti() {
      this.modalitaRicerca = true;
      await this.eseguiChiamataUtenti(this.testoRicerca.trim());
    },

    // Funzione centralizzata per GET /users conforme a api.yaml
    async eseguiChiamataUtenti(query) {
      this.caricamentoUtenti = true;
      try {
        // Se query è vuota chiama '/users', altrimenti '/users?username=query'
        const url = query ? `/users?username=${encodeURIComponent(query)}` : '/users';
        const res = await api.get(url);
        const lista = Array.isArray(res.data) ? res.data : [];

        // Filtra escludendo l'utente corrente autenticato
        this.utentiTrovati = lista.filter(u => {
          const nome = u.username || u.Username;
          return nome && nome.toLowerCase() !== this.mioUsername.toLowerCase();
        });
      } catch (e) {
        console.error("Errore nel recupero utenti:", e);
      } finally {
        this.caricamentoUtenti = false;
      }
    },

    // Annulla la ricerca e torna alle chat aperte
    annullaRicerca() {
      this.modalitaRicerca = false;
      this.testoRicerca = '';
      this.utentiTrovati = [];
    },

    // 4. Avvia o recupera una chat diretta (POST /conversations)
    async avviaChatDiretta(usernameSelezionato) {
      try {
        const res = await api.post('/conversations', { username: usernameSelezionato });
        this.annullaRicerca();
        await this.scaricaConversazioni();
        
        // Apre la conversazione appena creata o recuperata
        if (res.data && res.data.id) {
          this.apriChat(res.data.id);
        }
      } catch (e) {
        alert("Errore nell'avvio della conversazione diretta.");
        console.error(e);
      }
    },

    apriChat(id) {
      console.log("Apertura chat con ID:", id);
      // Fase successiva: recupero messaggi
    },

    gestisciGruppoCreato() {
      this.mostraModaleGruppo = false;
      this.scaricaConversazioni();
    },

    logout() {
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      this.$router.push('/login');
    }
  },
  mounted() {
    this.scaricaConversazioni();
  }
};
</script>

<style scoped>
.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>