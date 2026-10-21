<template>
  <div class="modal-overlay" @click.self="$emit('chiudi')">
    <div class="modal-content card p-4 shadow text-center">
      <!-- Intestazione del Modale -->
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h5 class="m-0 fw-bold">Il mio Profilo</h5>
        <button type="button" class="btn-close" @click="$emit('chiudi')"></button>
      </div>

      <!-- Avatar Foto Profilo (cliccabile per cambiare foto) -->
      <div class="my-2 d-flex flex-column align-items-center">
        <div class="position-relative profile-pic-wrapper mb-2" @click="apriSelettoreFile" title="Clicca per cambiare immagine">
          <img
            v-if="fotoAttuale"
            :src="fotoAttuale"
            class="rounded-circle shadow-sm border profile-img"
            alt="Foto Profilo"
          />
          <div
            v-else
            class="rounded-circle bg-primary text-white d-flex justify-content-center align-items-center shadow-sm profile-placeholder"
          >
            {{ nuovoUsername ? nuovoUsername.charAt(0).toUpperCase() : '?' }}
          </div>
          <div class="profile-pic-overlay rounded-circle d-flex justify-content-center align-items-center">
            <span class="small text-white fw-bold">Cambia</span>
          </div>
        </div>

        <!-- Input file nascosto per selezionare l'immagine da disco -->
        <input 
          type="file" 
          ref="fileInput" 
          @change="caricaImmagine" 
          accept="image/*" 
          class="d-none" 
        />
        <button 
          type="button" 
          class="btn btn-sm btn-outline-secondary" 
          @click="apriSelettoreFile"
          :disabled="inCaricamentoFoto"
        >
          {{ inCaricamentoFoto ? 'Caricamento...' : 'Seleziona nuova foto' }}
        </button>
      </div>

      <!-- Modifica Nome Utente (PUT /users/me/name) -->
      <div class="my-3 text-start">
        <label class="form-label small fw-bold text-dark mb-1">Modifica Username (3-16 caratteri):</label>
        <div class="input-group input-group-sm">
          <input 
            type="text" 
            v-model="nuovoUsername" 
            class="form-control" 
            minlength="3" 
            maxlength="16" 
            placeholder="Nuovo username..."
            :disabled="inSalvataggioNome"
          />
          <button 
            class="btn btn-primary" 
            @click="salvaNome" 
            :disabled="!puoSalvareNome || inSalvataggioNome"
          >
            {{ inSalvataggioNome ? 'Salvataggio...' : 'Salva' }}
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

      <!-- Azioni Logout e Chiusura -->
      <div class="d-flex flex-column gap-2 mt-2">
        <button type="button" @click="$emit('logout')" class="btn btn-outline-danger btn-sm w-100">
          Esci dall'Account
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
  name: 'ProfileModal',
  props: {
    username: {
      type: String,
      required: true
    },
    foto: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      nuovoUsername: this.username,
      fotoAttuale: this.foto,
      inCaricamentoFoto: false,
      inSalvataggioNome: false,
      messaggioErrore: '',
      messaggioSuccesso: ''
    };
  },
  computed: {
    puoSalvareNome() {
      const nomeTrim = this.nuovoUsername.trim();
      return (
        nomeTrim.length >= 3 &&
        nomeTrim.length <= 16 &&
        nomeTrim !== this.username
      );
    }
  },
  methods: {
    apriSelettoreFile() {
      this.$refs.fileInput.click();
    },

    // 1. Converte il file selezionato in Base64 e chiama PUT /users/me/propic
    async caricaImmagine(event) {
      const file = event.target.files[0];
      if (!file) return;

      this.messaggioErrore = '';
      this.messaggioSuccesso = '';
      this.inCaricamentoFoto = true;

      const reader = new FileReader();
      reader.onload = async () => {
        try {
          const base64Data = reader.result;
          await api.put('/users/me/propic', {
            picture: base64Data
          });

          this.fotoAttuale = base64Data;
          localStorage.setItem('proPic', base64Data);
          this.messaggioSuccesso = 'Foto profilo aggiornata con successo!';
          
          this.$emit('profiloAggiornato', {
            username: this.username,
            foto: base64Data
          });
        } catch (err) {
          console.error("Errore salvataggio foto profilo:", err);
          this.messaggioErrore = "Impossibile aggiornare la foto profilo.";
        } finally {
          this.inCaricamentoFoto = false;
          // Resetta il valore dell'input per permettere di ricaricare lo stesso file
          event.target.value = '';
        }
      };

      reader.onerror = () => {
        this.messaggioErrore = "Errore durante la lettura del file immagine.";
        this.inCaricamentoFoto = false;
      };

      reader.readAsDataURL(file);
    },

    // 2. Chiama PUT /users/me/name per aggiornare lo username
    async salvaNome() {
      const nomeValido = this.nuovoUsername.trim();
      if (!this.puoSalvareNome) return;

      this.messaggioErrore = '';
      this.messaggioSuccesso = '';
      this.inSalvataggioNome = true;

      try {
        await api.put('/users/me/name', {
          name: nomeValido
        });

        localStorage.setItem('username', nomeValido);
        this.messaggioSuccesso = 'Nome utente aggiornato!';

        this.$emit('profiloAggiornato', {
          username: nomeValido,
          foto: this.fotoAttuale
        });
      } catch (err) {
        console.error("Errore salvataggio username:", err);
        if (err.response && err.response.data && err.response.data.error) {
          this.messaggioErrore = err.response.data.error;
        } else {
          this.messaggioErrore = "Nome già in uso o non valido (3-16 caratteri).";
        }
      } finally {
        this.inSalvataggioNome = false;
      }
    }
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
  width: 380px;
  max-width: 95%;
  background-color: #ffffff !important;
  color: #212529 !important;
  border: none;
}

.profile-pic-wrapper {
  width: 96px;
  height: 96px;
  cursor: pointer;
}

.profile-img {
  width: 96px;
  height: 96px;
  object-fit: cover;
}

.profile-placeholder {
  width: 96px;
  height: 96px;
  font-size: 2.4rem;
  font-weight: bold;
}

.profile-pic-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.45);
  opacity: 0;
  transition: opacity 0.2s ease-in-out;
}

.profile-pic-wrapper:hover .profile-pic-overlay {
  opacity: 1;
}
</style>