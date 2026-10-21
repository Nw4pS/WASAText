<template>
  <div class="modal-overlay" @click.self="$emit('chiudi')">
    <div class="modal-content card p-4 shadow" role="dialog" aria-modal="true">
      <!-- Intestazione del Modale -->
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h5 class="m-0 fw-bold">Info Gruppo</h5>
        <button type="button" class="btn-close" @click="$emit('chiudi')" aria-label="Chiudi"></button>
      </div>

      <!-- Avatar Foto Gruppo -->
      <div class="my-2 d-flex flex-column align-items-center">
        <div 
          class="position-relative group-pic-wrapper mb-2" 
          @click="apriSelettoreFoto" 
          title="Clicca per cambiare l'immagine del gruppo"
        >
          <img
            v-if="fotoAttuale"
            :src="fotoAttuale"
            class="rounded-circle shadow-sm border group-img"
            alt="Foto Gruppo"
          />
          <div
            v-else
            class="rounded-circle bg-primary text-white d-flex justify-content-center align-items-center shadow-sm group-placeholder"
          >
            {{ nuovoNome ? nuovoNome.charAt(0).toUpperCase() : '?' }}
          </div>
          <div class="group-pic-overlay rounded-circle d-flex justify-content-center align-items-center">
            <span class="small text-white fw-bold">Cambia</span>
          </div>
        </div>

        <input 
          type="file" 
          ref="fileInputGruppo" 
          @change="caricaFotoGruppo" 
          accept="image/*" 
          class="d-none" 
        />
        <button 
          type="button" 
          class="btn btn-sm btn-outline-secondary" 
          @click="apriSelettoreFoto"
          :disabled="inCaricamentoFoto"
        >
          {{ inCaricamentoFoto ? 'Caricamento foto...' : 'Cambia foto gruppo' }}
        </button>
      </div>

      <!-- Modifica Nome del Gruppo (PUT /groups/{id}/name) -->
      <div class="my-3 text-start">
        <label class="form-label small fw-bold text-dark mb-1">Nome Gruppo:</label>
        <div class="input-group input-group-sm">
          <input 
            type="text" 
            v-model="nuovoNome" 
            class="form-control" 
            placeholder="Nome del gruppo..."
            :disabled="inSalvataggioNome"
            @keyup.enter="salvaNomeGruppo"
          />
          <button 
            class="btn btn-primary" 
            @click="salvaNomeGruppo" 
            :disabled="!puoSalvareNome || inSalvataggioNome"
          >
            {{ inSalvataggioNome ? 'Salvataggio...' : 'Salva Nome' }}
          </button>
        </div>
      </div>

      <!-- Feedback Utente -->
      <div v-if="messaggioErrore" class="alert alert-danger py-1 px-2 small mb-2 text-start">
        {{ messaggioErrore }}
      </div>
      <div v-if="messaggioSuccesso" class="alert alert-success py-1 px-2 small mb-2 text-start">
        {{ messaggioSuccesso }}
      </div>

      <!-- Elenco Partecipanti Attuali (GET /groups/{id}/members) -->
      <div class="mb-3 text-start">
        <label class="form-label small fw-bold text-dark mb-1">
          Partecipanti attuali ({{ membri.length }}):
        </label>
        <div class="border rounded p-2 bg-light members-box">
          <div v-if="caricamentoMembri" class="text-muted small p-1">
            Caricamento partecipanti in corso...
          </div>
          <div v-else-if="membri.length === 0" class="text-muted small p-1 fst-italic">
            Nessun partecipante trovato.
          </div>
          <div v-else class="d-flex flex-wrap gap-1">
            <span 
              v-for="m in membri" 
              :key="m" 
              class="badge bg-secondary text-white px-2 py-1"
            >
              {{ m }}
            </span>
          </div>
        </div>
      </div>

      <!-- Ricerca Contatti (con Debounce) -->
      <div class="mb-2 text-start">
        <label class="form-label small fw-bold text-dark mb-1">Cerca utenti da aggiungere:</label>
        <input 
          v-model="queryRicerca" 
          @input="gestisciInputRicerca" 
          type="text" 
          class="form-control form-control-sm" 
          placeholder="Filtra per username..." 
        />
      </div>

      <!-- Elenco Contatti Disponibili -->
      <div class="text-start mb-2">
        <label class="form-label small text-muted mb-1">Contatti disponibili (clicca per aggiungere):</label>
        <div class="border rounded p-2 bg-light utenti-box">
          <div v-if="caricamentoUtenti" class="text-muted small p-1">
            Caricamento contatti in corso...
          </div>
          <div v-else-if="contattiFiltrati.length === 0" class="text-muted small p-1 fst-italic">
            Nessun nuovo contatto disponibile.
          </div>
          
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
      </div>

      <!-- Membri Selezionati per l'aggiunta -->
      <div class="mb-3 text-start">
        <label class="form-label small fw-bold text-dark mb-1">
          Nuovi utenti da aggiungere ({{ membriDaAggiungere.length }}):
        </label>
        <div>
          <span v-if="membriDaAggiungere.length === 0" class="text-muted small fst-italic">
            Nessun nuovo utente ancora selezionato.
          </span>
          <span 
            v-for="(m, i) in membriDaAggiungere" 
            :key="m" 
            class="badge bg-primary me-1 p-2"
          >
            {{ m }}
            <span role="button" @click="rimuoviDaAggiungere(i)" class="ms-1 fw-bold">&times;</span>
          </span>
        </div>
      </div>

      <!-- Pulsanti Azione -->
      <div class="d-flex flex-column gap-2 mt-2">
        <button 
          v-if="membriDaAggiungere.length > 0"
          type="button" 
          @click="confermaAggiunta" 
          class="btn btn-success btn-sm w-100"
          :disabled="inInvioMembri"
        >
          {{ inInvioMembri ? 'Aggiunta in corso...' : 'Salva Nuovi Membri nel Gruppo' }}
        </button>

        <button type="button" @click="$emit('abbandona')" class="btn btn-outline-danger btn-sm w-100">
          Abbandona Gruppo
        </button>
        <button type="button" @click="$emit('chiudi')" class="btn btn-secondary btn-sm w-100">
          Chiudi
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/axios.js';

export default {
  name: 'GroupInfoModal',
  props: {
    gruppo: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      nuovoNome: '',
      fotoAttuale: '',
      membri: [],
      tuttiGliUtenti: [],
      membriDaAggiungere: [],
      queryRicerca: '',
      caricamentoMembri: false,
      caricamentoUtenti: false,
      inCaricamentoFoto: false,
      inSalvataggioNome: false,
      inInvioMembri: false,
      messaggioErrore: '',
      messaggioSuccesso: '',
      timerRicerca: null,
      timerSuccesso: null
    };
  },
  watch: {
    // Garantisce che il form si sincronizzi sempre se cambia l'oggetto prop
    gruppo: {
      immediate: true,
      handler(nuovoValore) {
        if (nuovoValore) {
          this.nuovoNome = nuovoValore.name || '';
          this.fotoAttuale = nuovoValore.picture || '';
        }
      }
    }
  },
  computed: {
    puoSalvareNome() {
      const nomeTrim = this.nuovoNome.trim();
      return nomeTrim.length > 0 && nomeTrim !== this.gruppo.name;
    },
    contattiFiltrati() {
      return this.tuttiGliUtenti.filter(u => {
        const username = u.username || u.Username;
        if (!username) return false;
        const giaNelGruppo = this.membri.includes(username);
        const giaSelezionato = this.membriDaAggiungere.includes(username);
        return !giaNelGruppo && !giaSelezionato;
      });
    }
  },
  methods: {
    apriSelettoreFoto() {
      this.$refs.fileInputGruppo.click();
    },

    mostraSuccessoTemporaneo(messaggio) {
      this.messaggioSuccesso = messaggio;
      if (this.timerSuccesso) clearTimeout(this.timerSuccesso);
      this.timerSuccesso = setTimeout(() => {
        this.messaggioSuccesso = '';
      }, 3000);
    },

    // 1. Modifica Foto Gruppo (PUT /groups/{id}/photo)
    async caricaFotoGruppo(event) {
      const file = event.target.files[0];
      if (!file) return;

      this.messaggioErrore = '';
      this.inCaricamentoFoto = true;

      const reader = new FileReader();
      reader.onload = async () => {
        try {
          const base64Data = reader.result;
          await api.put(`/groups/${this.gruppo.id}/photo`, {
            picture: base64Data
          });

          this.fotoAttuale = base64Data;
          this.mostraSuccessoTemporaneo('Foto del gruppo aggiornata!');

          this.$emit('gruppoAggiornato', {
            id: this.gruppo.id,
            name: this.nuovoNome,
            picture: base64Data
          });
        } catch (err) {
          console.error("Errore aggiornamento foto gruppo:", err);
          this.messaggioErrore = "Impossibile aggiornare la foto del gruppo.";
        } finally {
          this.inCaricamentoFoto = false;
          event.target.value = '';
        }
      };

      reader.onerror = () => {
        this.messaggioErrore = "Errore durante la lettura del file immagine.";
        this.inCaricamentoFoto = false;
      };

      reader.readAsDataURL(file);
    },

    // 2. Modifica Nome Gruppo (PUT /groups/{id}/name)
    async salvaNomeGruppo() {
      const nomeValido = this.nuovoNome.trim();
      if (!this.puoSalvareNome) return;

      this.messaggioErrore = '';
      this.inSalvataggioNome = true;

      try {
        await api.put(`/groups/${this.gruppo.id}/name`, {
          name: nomeValido
        });

        this.mostraSuccessoTemporaneo('Nome del gruppo aggiornato!');

        this.$emit('gruppoAggiornato', {
          id: this.gruppo.id,
          name: nomeValido,
          picture: this.fotoAttuale
        });
      } catch (err) {
        console.error("Errore aggiornamento nome gruppo:", err);
        if (err.response && err.response.data && err.response.data.error) {
          this.messaggioErrore = err.response.data.error;
        } else {
          this.messaggioErrore = "Impossibile aggiornare il nome del gruppo.";
        }
      } finally {
        this.inSalvataggioNome = false;
      }
    },

    // 3. Recupero Partecipanti Attuali (GET /groups/{id}/members)
    async caricaMembri() {
      if (!this.gruppo || !this.gruppo.id) return;
      this.caricamentoMembri = true;
      try {
        const res = await api.get(`/groups/${this.gruppo.id}/members`);
        this.membri = Array.isArray(res.data) ? res.data : [];
      } catch (err) {
        console.error("Errore nel recupero partecipanti del gruppo:", err);
        this.membri = [];
      } finally {
        this.caricamentoMembri = false;
      }
    },

    // 4. Ricerca Contatti con Debounce anti-race condition
    gestisciInputRicerca() {
      if (this.timerRicerca) clearTimeout(this.timerRicerca);
      this.timerRicerca = setTimeout(() => {
        this.eseguiRicercaUtenti();
      }, 300);
    },

    async eseguiRicercaUtenti() {
      this.caricamentoUtenti = true;
      try {
        const q = this.queryRicerca.trim();
        const url = q ? `/users?username=${encodeURIComponent(q)}` : '/users';
        const res = await api.get(url);
        this.tuttiGliUtenti = Array.isArray(res.data) ? res.data : [];
      } catch (err) {
        console.error("Errore nel recupero utenti:", err);
      } finally {
        this.caricamentoUtenti = false;
      }
    },

    selezionaUtente(username) {
      if (!this.membriDaAggiungere.includes(username)) {
        this.membriDaAggiungere.push(username);
      }
    },

    rimuoviDaAggiungere(index) {
      this.membriDaAggiungere.splice(index, 1);
    },

    // 5. Aggiunta Membri (POST /groups/{id}/members)
    async confermaAggiunta() {
      if (this.membriDaAggiungere.length === 0) return;
      this.inInvioMembri = true;
      this.messaggioErrore = '';

      try {
        await api.post(`/groups/${this.gruppo.id}/members`, {
          usernames: this.membriDaAggiungere
        });
        this.membriDaAggiungere = [];
        this.mostraSuccessoTemporaneo('Nuovi partecipanti aggiunti con successo!');
        await this.caricaMembri();
        this.$emit('membroAggiunto');
      } catch (err) {
        console.error("Errore durante l'aggiunta membri:", err);
        if (err.response && err.response.data && err.response.data.error) {
          this.messaggioErrore = err.response.data.error;
        } else {
          this.messaggioErrore = "Impossibile aggiungere gli utenti selezionati.";
        }
      } finally {
        this.inInvioMembri = false;
      }
    },

    gestisciKeyDown(event) {
      if (event.key === 'Escape') {
        this.$emit('chiudi');
      }
    }
  },
  mounted() {
    this.caricaMembri();
    this.eseguiRicercaUtenti();
    window.addEventListener('keydown', this.gestisciKeyDown);
  },
  beforeUnmount() {
    window.removeEventListener('keydown', this.gestisciKeyDown);
    if (this.timerRicerca) clearTimeout(this.timerRicerca);
    if (this.timerSuccesso) clearTimeout(this.timerSuccesso);
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
  width: 440px;
  max-width: 95%;
  background-color: #ffffff !important;
  color: #212529 !important;
  border: none;
  max-height: 90vh;
  overflow-y: auto;
}

.group-pic-wrapper {
  width: 90px;
  height: 90px;
  cursor: pointer;
}

.group-img {
  width: 90px;
  height: 90px;
  object-fit: cover;
}

.group-placeholder {
  width: 90px;
  height: 90px;
  font-size: 2.2rem;
  font-weight: bold;
}

.group-pic-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.45);
  opacity: 0;
  transition: opacity 0.2s ease-in-out;
}

.group-pic-wrapper:hover .group-pic-overlay {
  opacity: 1;
}

.members-box {
  min-height: 48px;
  max-height: 100px;
  overflow-y: auto;
  background-color: #f8f9fa !important;
  border: 1px solid #ced4da !important;
}

.utenti-box {
  min-height: 52px;
  max-height: 120px;
  overflow-y: auto;
  background-color: #f8f9fa !important;
  border: 1px solid #ced4da !important;
}
</style>