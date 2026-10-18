<template>
  <div class="modal-overlay">
    <div class="modal-content card p-4 shadow">
      <h5 class="mb-3">Crea un nuovo gruppo</h5>
      
      <!-- Nome del gruppo -->
      <div class="mb-3">
        <label class="form-label fw-bold">Nome del gruppo</label>
        <input 
          v-model="nomeGruppo" 
          type="text" 
          class="form-control" 
          placeholder="Es. Fantacalcio" 
        />
      </div>

      <!-- Ricerca contatti -->
      <div class="mb-3">
        <label class="form-label fw-bold">Cerca o seleziona utenti da aggiungere</label>
        <input 
          v-model="queryRicerca" 
          @input="eseguiRicerca" 
          type="text" 
          class="form-control" 
          placeholder="Filtra per username..." 
        />
      </div>

      <!-- Box con l'elenco dei contatti disponibili -->
      <label class="form-label small text-muted">Contatti disponibili (clicca per aggiungere):</label>
      <div class="border rounded p-2 bg-light mb-3 utenti-box">
        <div v-if="caricamento" class="text-muted small p-1">
          Caricamento contatti in corso...
        </div>
        <div v-else-if="contattiFiltrati.length === 0" class="text-muted small p-1 fst-italic">
          Nessun contatto disponibile o corrispondente alla ricerca.
        </div>
        
        <!-- Pulsanti dei contatti -->
        <button
          v-for="u in contattiFiltrati"
          :key="u.username"
          type="button"
          @click="selezionaUtente(u.username)"
          class="btn btn-sm btn-outline-primary m-1"
        >
          + {{ u.username }}
        </button>
      </div>

      <!-- Membri selezionati -->
      <div class="mb-3">
        <label class="form-label small fw-bold">Membri selezionati ({{ membriSelezionati.length }}):</label>
        <div>
          <span v-if="membriSelezionati.length === 0" class="text-muted small fst-italic">
            Nessun membro ancora selezionato.
          </span>
          <span 
            v-for="(m, i) in membriSelezionati" 
            :key="m" 
            class="badge bg-primary me-1 p-2"
          >
            {{ m }}
            <span role="button" @click="rimuoviMembro(i)" class="ms-1 fw-bold">&times;</span>
          </span>
        </div>
      </div>

      <!-- Pulsanti Azione -->
      <div class="d-flex justify-content-end gap-2 mt-2">
        <button type="button" @click="$emit('chiudi')" class="btn btn-secondary btn-sm">
          Annulla
        </button>
        <button 
          type="button" 
          @click="confermaCreazione" 
          class="btn btn-success btn-sm"
          :disabled="!nomeGruppo.trim() || membriSelezionati.length === 0"
        >
          Crea Gruppo
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/axios.js';

export default {
  name: 'NewGroupModal',
  props: {
    mioUsername: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      nomeGruppo: '',
      queryRicerca: '',
      tuttiGliUtenti: [],
      membriSelezionati: [],
      caricamento: false
    };
  },
  computed: {
    utenteLoggato() {
      return this.mioUsername || localStorage.getItem('username') || '';
    },
    contattiFiltrati() {
      return this.tuttiGliUtenti.filter(u => {
        const username = u.username || u.Username;
        if (!username) return false;
        const nonEMeStesso = username.toLowerCase() !== this.utenteLoggato.toLowerCase();
        const nonAncoraScelto = !this.membriSelezionati.includes(username);
        return nonEMeStesso && nonAncoraScelto;
      });
    }
  },
  methods: {
    // Chiama GET /users (con o senza parametro username)
    async eseguiRicerca() {
      this.caricamento = true;
      try {
        const q = this.queryRicerca.trim();
        const url = q ? `/users?username=${encodeURIComponent(q)}` : '/users';
        const res = await api.get(url);
        this.tuttiGliUtenti = Array.isArray(res.data) ? res.data : [];
      } catch (err) {
        console.error("Errore nel recupero utenti per gruppo:", err);
      } finally {
        this.caricamento = false;
      }
    },
    selezionaUtente(username) {
      if (!this.membriSelezionati.includes(username)) {
        this.membriSelezionati.push(username);
      }
    },
    rimuoviMembro(indice) {
      this.membriSelezionati.splice(indice, 1);
    },
    async confermaCreazione() {
      try {
        await api.post('/groups', {
          name: this.nomeGruppo.trim(),
          members: this.membriSelezionati
        });
        this.$emit('fatto');
      } catch (err) {
        alert("Errore durante la creazione del gruppo.");
        console.error(err);
      }
    }
  },
  mounted() {
    // Carica immediatamente tutti gli utenti appena il modale viene aperto
    this.eseguiRicerca();
  }
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: rgba(0, 0, 0, 0.55);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1050;
}
.modal-content {
  width: 450px;
  max-width: 95%;
}
.utenti-box {
  min-height: 52px;
  max-height: 140px;
  overflow-y: auto;
}
</style>