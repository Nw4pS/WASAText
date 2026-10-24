<template>
  <div class="container-fluid vh-100 d-flex flex-column p-0">
    <!-- Header superiore -->
    <header class="bg-primary text-white p-2 px-3 d-flex justify-content-between align-items-center shadow-sm">
      <h5 class="m-0 fw-bold">WASAText</h5>
      
      <!-- Pulsante Profilo con Foto e Nome Utente -->
      <button 
        @click="mostraModaleProfilo = true" 
        class="btn btn-link text-white text-decoration-none d-flex align-items-center gap-2 p-1 rounded hover-profile"
        title="Visualizza profilo"
      >
        <img 
          v-if="miaFoto" 
          :src="miaFoto" 
          class="rounded-circle border border-white" 
          style="width: 36px; height: 36px; object-fit: cover;" 
          alt="Foto Profilo"
        />
        <div 
          v-else 
          class="rounded-circle bg-white text-primary d-flex justify-content-center align-items-center fw-bold shadow-sm" 
          style="width: 36px; height: 36px; font-size: 1rem;"
        >
          {{ mioUsername ? mioUsername.charAt(0).toUpperCase() : '?' }}
        </div>
        <span class="fw-semibold">{{ mioUsername }}</span>
      </button>
    </header>

    <div class="row flex-grow-1 g-0 overflow-hidden">
      <!-- Sidebar sinistra: 4 colonne -->
      <aside class="col-4 border-end bg-light d-flex flex-column h-100">
        
        <!-- Barra di ricerca e pulsante Nuovo Gruppo affiancati -->
        <div class="p-3 border-bottom bg-white">
          <div class="d-flex gap-2">
            <div class="input-group flex-grow-1">
              <input 
                type="text" 
                class="form-control" 
                placeholder="Cerca utente per iniziare chat..." 
                v-model="testoRicerca"
                @focus="attivaModalitaRicerca"
                @input="cercaUtenti"
              />
              <button 
                v-if="modalitaRicerca" 
                @click="annullaRicerca" 
                class="btn btn-outline-secondary" 
                type="button"
              >
                &times;
              </button>
            </div>

            <button 
              @click="mostraModaleGruppo = true" 
              class="btn btn-outline-primary d-flex align-items-center" 
              title="Crea nuovo gruppo"
            >
              + Gruppo
            </button>
          </div>
        </div>

        <div class="list-group list-group-flush overflow-auto flex-grow-1">
          
          <!-- SEZIONE RICERCA UTENTI -->
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
                <div class="small text-muted">Clicca per aprire la chat</div>
              </div>
            </button>
          </div>

          <!-- SEZIONE LISTA CHAT ESISTENTI -->
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
              :class="{ 'active-conversation': conv.id === idChatAttiva }"
              @seleziona="apriChat(conv)" 
            />

            <div v-if="conversazioni.length === 0 && !caricamentoChat" class="p-4 text-center text-muted">
              Nessuna chat attiva. Usa la barra di ricerca in alto per iniziare a chattare!
            </div>
          </div>

        </div>
      </aside>

      <!-- Area principale destra: 8 colonne -->
      <main class="col-8 d-flex flex-column h-100 bg-white">
        
        <!-- Placeholder visibile quando nessuna chat è selezionata -->
        <div v-if="!idChatAttiva" class="d-flex flex-column justify-content-center align-items-center h-100 text-muted">
          <div class="p-4 text-center">
            <h4 class="fw-bold">Nessuna conversazione selezionata</h4>
            <p class="mb-0">Seleziona una conversazione dall'elenco a sinistra o cerca un utente.</p>
          </div>
        </div>

        <!-- Finestra della chat attiva -->
        <div v-else class="d-flex flex-column h-100">
          
          <!-- Header della chat attiva: foto e titolo resi interattivi se è un gruppo -->
        <div class="p-3 border-bottom d-flex align-items-center justify-content-between bg-light">
          <div 
            class="d-flex align-items-center"
            :class="{ 'group-clickable p-1 rounded': chatAttivaInfo && chatAttivaInfo.isGroup }"
            @click="chatAttivaInfo && chatAttivaInfo.isGroup ? (mostraModaleInfoGruppo = true) : null"
            :title="chatAttivaInfo && chatAttivaInfo.isGroup ? 'Clicca per dettagli gruppo' : ''"
          >
            <!-- Foto del gruppo o della chat -->
            <img 
              v-if="chatAttivaInfo && chatAttivaInfo.picture" 
              :src="chatAttivaInfo.picture" 
              class="rounded-circle border me-3" 
              style="width: 42px; height: 42px; object-fit: cover;" 
              alt="Foto Conversazione"
            />
            <div 
              v-else 
              class="rounded-circle bg-primary text-white d-flex justify-content-center align-items-center me-3 shadow-sm" 
              style="width: 42px; height: 42px; font-weight: bold;"
            >
              {{ chatAttivaInfo ? chatAttivaInfo.name.charAt(0).toUpperCase() : '?' }}
            </div>

            <div>
              <h6 class="mb-0 fw-bold d-flex align-items-center gap-2">
                {{ chatAttivaInfo ? chatAttivaInfo.name : 'Chat' }}
                <span 
                  v-if="chatAttivaInfo && chatAttivaInfo.isGroup" 
                  class="badge bg-primary-subtle text-primary border border-primary-subtle" 
                  style="font-size: 0.65rem;"
                >
                  Info
                </span>
              </h6>
              <small class="text-muted">
                {{ chatAttivaInfo && chatAttivaInfo.isGroup ? 'Gruppo (clicca per info)' : 'Chat diretta' }}
              </small>
            </div>
          </div>

          <!-- Pulsante rapido opzionale -->
          <button 
            v-if="chatAttivaInfo && chatAttivaInfo.isGroup" 
            @click="abbandonaGruppo" 
            class="btn btn-outline-danger btn-sm"
          >
            Abbandona Gruppo
          </button>
        </div>

          <!-- Area messaggi con scroll -->
          <div class="flex-grow-1 overflow-auto p-3 bg-chat" ref="areaMessaggi" @click="chiudiMenuMessaggio">
            <div v-if="caricamentoMessaggi" class="text-center text-muted mt-4">
              Caricamento messaggi in corso...
            </div>
            
            <div v-else-if="messaggi.length === 0" class="text-center text-muted mt-4">
              Nessun messaggio presente. Scrivi il primo messaggio!
            </div>

            <!-- Lista dei messaggi (bolle stile chat) -->
            <div 
              v-else
              v-for="msg in messaggi" 
              :key="msg.id" 
              class="mb-3 d-flex"
              :class="msg.user === mioUsername ? 'justify-content-end' : 'justify-content-start'"
            >
              <div 
                class="p-2 px-3 rounded shadow-sm message-bubble position-relative"
                :class="msg.user === mioUsername ? 'bg-primary text-white' : 'bg-white border text-dark'"
              >
                <!-- Header bolla: Mittente + Tasto Tre Puntini (⋮) -->
                <div class="d-flex justify-content-between align-items-center mb-1">
                  <span v-if="msg.user !== mioUsername" class="small fw-bold text-primary">
                    {{ msg.user }}
                  </span>
                  <span v-else></span>

                  <!-- Pulsante Tre Puntini Verticali -->
                  <button 
                    type="button" 
                    class="btn btn-sm btn-link p-0 text-decoration-none dots-btn"
                    :class="msg.user === mioUsername ? 'text-white' : 'text-muted'"
                    @click.stop="toggleMenuMessaggio(msg.id)"
                    title="Opzioni messaggio"
                  >
                    &#8942;
                  </button>

                  <!-- Menu a comparsa (Dropdown) -->
                  <div 
                    v-if="menuApertoId === msg.id" 
                    class="context-menu shadow border rounded bg-white py-1 text-dark"
                    :class="msg.user === mioUsername ? 'menu-destra' : 'menu-sinistra'"
                    @click.stop
                  >
                    <!-- 1. Emoji rapide (PUT /messages/{id}/reactions) -->
                    <div class="px-2 py-1 border-bottom">
                      <span class="small text-muted d-block" style="font-size: 0.7rem;">Reagisci:</span>
                      <div class="d-flex gap-1">
                        <button 
                          v-for="emoji in ['👍', '❤️', '😂', '😮', '😢', '🔥']" 
                          :key="emoji"
                          type="button" 
                          class="btn btn-sm btn-light p-1 border-0 emoji-btn"
                          @click="aggiungiReazione(msg.id, emoji)"
                        >
                          {{ emoji }}
                        </button>
                      </div>
                    </div>

                    <!-- 2. Rimuovi reazione propria (DELETE /messages/{id}/reactions) -->
                    <button 
                      v-if="haMiaReazione(msg)" 
                      type="button" 
                      class="dropdown-item small text-danger py-1"
                      @click="rimuoviReazione(msg.id)"
                    >
                      ❌ Rimuovi reazione
                    </button>

                    <!-- 3. Rispondi -->
                    <button type="button" class="dropdown-item small py-1" @click="impostaRisposta(msg)">
                      ↩ Rispondi
                    </button>

                    <!-- 4. Inoltra -->
                    <button type="button" class="dropdown-item small py-1" @click="apriInoltro(msg)">
                      ↪ Inoltra
                    </button>

                    <!-- 5. Elimina messaggio (solo se inviato da me) -->
                    <button 
                      v-if="msg.user === mioUsername" 
                      type="button" 
                      class="dropdown-item small text-danger border-top mt-1 py-1"
                      @click="eliminaMessaggio(msg.id)"
                    >
                      🗑 Elimina messaggio
                    </button>
                  </div>
                </div>

                <!-- Badge Inoltrato -->
                <div v-if="msg.isForwarded" class="small fst-italic text-warning mb-1">
                  ↪ Inoltrato
                </div>

                <!-- se il messaggio risponde a un altro (answerTo) -->
                <!-- Anteprima ricca del messaggio citato (risposta) -->
                <div 
                  v-if="msg.answerTo" 
                  class="p-1 px-2 mb-1 rounded small border-start border-3 text-start"
                  :class="msg.user === mioUsername ? 'bg-white bg-opacity-25 border-white text-white' : 'bg-light border-primary text-muted'"
                  style="font-size: 0.75rem; max-width: 100%;"
                >
                  <template v-if="recuperaMessaggioOriginale(msg.answerTo)">
                    <!-- Autore del messaggio a cui si risponde -->
                    <div class="fw-bold" style="font-size: 0.72rem;">
                      {{ recuperaMessaggioOriginale(msg.answerTo).user === mioUsername ? 'Tu' : recuperaMessaggioOriginale(msg.answerTo).user }}
                    </div>
                    <!-- Contenuto del messaggio citato (testo o etichetta foto) -->
                    <div class="text-truncate fst-italic" style="max-height: 2.4em;">
                      <span v-if="recuperaMessaggioOriginale(msg.answerTo).isImage">📷 Immagine</span>
                      <span v-else>{{ recuperaMessaggioOriginale(msg.answerTo).content }}</span>
                    </div>
                  </template>
                  
                  <!-- Fallback se il messaggio originale è stato cancellato o non è ancora in memoria -->
                  <template v-else>
                    <span class="fst-italic">In risposta al messaggio #{{ msg.answerTo }}</span>
                  </template>
                </div>

                <!-- Contenuto: immagine o testo -->
                <div v-if="msg.isImage">
                  <img :src="msg.content" class="img-fluid rounded my-1" style="max-height: 250px;" alt="Immagine allegata" />
                </div>
                <div v-else class="text-break" style="white-space: pre-wrap;">{{ msg.content }}</div>

                <!-- Footer messaggio: orario + stato (checkmarks) -->
                <div class="d-flex justify-content-end align-items-center gap-1 mt-1 small" style="font-size: 0.72rem; opacity: 0.85;">
                  <span>{{ formattaData(msg.timestamp) }}</span>
                  
                  <!-- Checkmarks conformi al PDF WASAText (per i propri messaggi inviati) -->
                  <span v-if="msg.user === mioUsername" class="ms-1 fw-bold">
                    <span v-if="msg.state === 'read'" title="Letto" class="text-info">✓✓</span>
                    <span v-else-if="msg.state === 'received'" title="Ricevuto">✓</span>
                    <span v-else title="Inviato">✓</span>
                  </span>
                </div>

                <!-- Reazioni / Emoji con username (specifica WASAText) -->
                <div v-if="msg.reactions && msg.reactions.length > 0" class="d-flex flex-wrap gap-1 mt-1 pt-1 border-top border-light-subtle">
                  <span 
                    v-for="react in msg.reactions" 
                    :key="react.id" 
                    class="badge bg-light text-dark border small" 
                    :title="'Aggiunto da ' + react.user"
                  >
                    {{ react.content }} <small class="text-muted">({{ react.user }})</small>
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Barra di digitazione messaggio -->
          <div class="border-top bg-light">
            <!-- Box anteprima risposta attiva -->
            <div 
              v-if="messaggioRisposta" 
              class="p-2 px-3 bg-white border-bottom d-flex justify-content-between align-items-center"
            >
              <div class="small text-truncate">
                <span class="fw-bold text-primary">Rispondendo a {{ messaggioRisposta.user }}: </span>
                <span class="text-muted fst-italic">{{ messaggioRisposta.content }}</span>
              </div>
              <button type="button" class="btn btn-sm btn-link text-muted p-0 ms-2" @click="messaggioRisposta = null">
                &times;
              </button>
            </div>

            <!-- Anteprima immagine selezionata -->
            <div 
              v-if="immagineSelezionata" 
              class="p-2 px-3 bg-white border-bottom d-flex align-items-center justify-content-between"
            >
              <div class="d-flex align-items-center gap-2">
                <img :src="immagineSelezionata" class="rounded border" style="width: 48px; height: 48px; object-fit: cover;" />
                <span class="small text-muted">{{ nomeFileImmagine }}</span>
              </div>
              <button type="button" class="btn btn-sm btn-link text-danger text-decoration-none" @click="annullaImmagine">
                &times; Rimuovi
              </button>
            </div>

            <div class="p-3">
              <form @submit.prevent="inviaMessaggio" class="input-group">
                <input 
                  type="file" 
                  ref="fileInputFoto" 
                  @change="selezionaFoto" 
                  accept="image/*" 
                  class="d-none" 
                />

                <button 
                  type="button" 
                  class="btn btn-outline-secondary" 
                  @click="$refs.fileInputFoto.click()"
                  title="Invia immagine"
                >
                  📷
                </button>

                <input 
                  type="text" 
                  class="form-control" 
                  :placeholder="immagineSelezionata ? 'Premi Invia per spedire l\'immagine...' : 'Scrivi un messaggio...'" 
                  v-model="nuovoMessaggio"
                  :disabled="caricamentoMessaggi || !!immagineSelezionata"
                  autofocus
                />

                <button 
                  class="btn btn-primary px-4" 
                  type="submit" 
                  :disabled="(!nuovoMessaggio.trim() && !immagineSelezionata) || caricamentoMessaggi"
                >
                  Invia
                </button>
              </form>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- Modale Profilo Utente -->
    <ProfileModal
      v-if="mostraModaleProfilo"
      :username="mioUsername"
      :foto="miaFoto"
      @chiudi="mostraModaleProfilo = false"
      @logout="logout"
      @profiloAggiornato="gestisciProfiloAggiornato"
    />

    <!-- Modale creazione gruppo (UNICO) -->
    <NewGroupModal 
      v-if="mostraModaleGruppo" 
      :mioUsername="mioUsername" 
      @chiudi="mostraModaleGruppo = false" 
      @fatto="gestisciGruppoCreato" 
    />
    <!-- modale dettagli gruppo -->
    <GroupInfoModal
      v-if="mostraModaleInfoGruppo && chatAttivaInfo && chatAttivaInfo.isGroup"
      :gruppo="chatAttivaInfo"
      @chiudi="mostraModaleInfoGruppo = false"
      @abbandona="gestisciAbbandonoDaModale"
      @membroAggiunto="scaricaConversazioni"
      @gruppoAggiornato="gestisciGruppoAggiornato"
    />
    <!-- Modale per Inoltrare il Messaggio -->
    <div v-if="mostraModaleInoltro" class="modal-overlay" @click.self="mostraModaleInoltro = false">
      <div class="modal-content card p-4 shadow">
        <div class="d-flex justify-content-between align-items-center mb-3">
          <h5 class="m-0 fw-bold">Inoltra messaggio</h5>
          <button type="button" class="btn-close" @click="mostraModaleInoltro = false"></button>
        </div>
        <p class="small text-muted mb-2">Seleziona la conversazione di destinazione:</p>
        <div class="list-group overflow-auto mb-3" style="max-height: 220px;">
          <button 
            v-for="c in conversazioni" 
            :key="c.id" 
            type="button" 
            class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
            @click="confermaInoltro(c.id)"
          >
            <span>{{ c.name }} <small class="text-muted">({{ c.isGroup ? 'Gruppo' : 'Diretta' }})</small></span>
            <span class="btn btn-sm btn-outline-primary py-0">Invia</span>
          </button>
        </div>
        <button type="button" class="btn btn-secondary btn-sm w-100" @click="mostraModaleInoltro = false">
          Annulla
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../services/axios.js';
import ConversationItem from '../components/ConversationItem.vue';
import NewGroupModal from '../components/NewGroupModal.vue';
import ProfileModal from '../components/ProfileModal.vue';
import GroupInfoModal from '../components/GroupInfoModal.vue';

export default {
  name: 'HomeView',
  components: {
    ConversationItem,
    NewGroupModal,
    ProfileModal,
    GroupInfoModal
  },
  data() {
    return {
      mioUsername: localStorage.getItem('username') || '',
      miaFoto: localStorage.getItem('proPic') || '',
      mostraModaleProfilo: false,
      mostraModaleGruppo: false,
      mostraModaleInfoGruppo: false,
      conversazioni: [],
      caricamentoChat: false,
      modalitaRicerca: false,
      testoRicerca: '',
      utentiTrovati: [],
      caricamentoUtenti: false,
      idChatAttiva: null,
      chatAttivaInfo: null,
      messaggi: [],
      nuovoMessaggio: '',
      caricamentoMessaggi: false,
      menuApertoId: null,           // Memorizza quale messaggio ha il menu ⋮ aperto
      messaggioRisposta: null,      // Oggetto messaggio a cui si risponde
      mostraModaleInoltro: false,   // Visibilità modale inoltro
      messaggioDaInoltrareId: null,  // ID del messaggio da inoltrare
      immagineSelezionata: null,
      nomeFileImmagine: ''
    };
  },
  methods: {
    recuperaMessaggioOriginale(idMessaggio) {
      if (!idMessaggio) return null;
      return this.messaggi.find(m => m.id === idMessaggio) || null;
    },   
  // Apertura/Chiusura menu messaggio⋮
    toggleMenuMessaggio(id) {
      this.menuApertoId = this.menuApertoId === id ? null : id;
    },
    chiudiMenuMessaggio() {
      this.menuApertoId = null;
    },

    // Rispondi
    impostaRisposta(msg) {
      this.messaggioRisposta = msg;
      this.chiudiMenuMessaggio();
    },

    // Aggiungi Reazione (PUT /messages/{id}/reactions - commentMessage)
    async aggiungiReazione(messaggioId, emoji) {
      this.chiudiMenuMessaggio();
      try {
        await api.put(`/messages/${messaggioId}/reactions`, { emoji });
        await this.apriChat(this.idChatAttiva);
      } catch (err) {
        console.error("Errore aggiunta reazione:", err);
      }
    },

    // Verifica se l'utente ha già reagito
    haMiaReazione(msg) {
      if (!msg.reactions || !Array.isArray(msg.reactions)) return false;
      return msg.reactions.some(r => r.user === this.mioUsername);
    },

    // Rimuovi Reazione (DELETE /messages/{id}/reactions - uncommentMessage)
    async rimuoviReazione(messaggioId) {
      this.chiudiMenuMessaggio();
      try {
        await api.delete(`/messages/${messaggioId}/reactions`);
        await this.apriChat(this.idChatAttiva);
      } catch (err) {
        console.error("Errore rimozione reazione:", err);
      }
    },

    // Inoltro messaggio (POST /conversations/{id}/messages/forward - forwardMessage)
    apriInoltro(msg) {
      this.chiudiMenuMessaggio();
      this.messaggioDaInoltrareId = msg.id;
      this.mostraModaleInoltro = true;
    },
    async confermaInoltro(targetConvId) {
      if (!this.messaggioDaInoltrareId) return;
      try {
        await api.post(`/conversations/${targetConvId}/messages/forward`, {
          messageId: this.messaggioDaInoltrareId
        });
        alert("Messaggio inoltrato con successo!");
        this.mostraModaleInoltro = false;
        this.messaggioDaInoltrareId = null;
        if (targetConvId === this.idChatAttiva) {
          await this.apriChat(this.idChatAttiva);
        }
        this.scaricaConversazioni();
      } catch (err) {
        console.error("Errore inoltro messaggio:", err);
        alert("Impossibile inoltrare il messaggio.");
      }
    },

    // Elimina messaggio (DELETE /messages/{id} - deleteMessage)
    async eliminaMessaggio(messaggioId) {
      this.chiudiMenuMessaggio();
      const conferma = confirm("Sei sicuro di voler eliminare questo messaggio?");
      if (!conferma) return;

      try {
        await api.delete(`/messages/${messaggioId}`);
        this.messaggi = this.messaggi.filter(m => m.id !== messaggioId);
        this.scaricaConversazioni();
      } catch (err) {
        console.error("Errore eliminazione messaggio:", err);
        alert("Impossibile eliminare il messaggio.");
      }
    },
    // 1. Scarica la lista delle conversazioni (GET /conversations)
    async scaricaConversazioni() {
      this.caricamentoChat = true;
      try {
        const risposta = await api.get('/conversations');
        let chatRicevute = risposta.data || [];
        
        chatRicevute.sort((a, b) => {
          if (!a.lastMessageDate) return 1;
          if (!b.lastMessageDate) return -1;
          return new Date(b.lastMessageDate) - new Date(a.lastMessageDate);
        });
        
        this.conversazioni = chatRicevute;

        // Se una chat era già aperta, aggiorniamo i suoi metadati
        if (this.idChatAttiva) {
          this.chatAttivaInfo = this.conversazioni.find(c => c.id === this.idChatAttiva) || this.chatAttivaInfo;
        }
      } catch (e) {
        console.error("Errore durante il download delle conversazioni", e);
      } finally {
        this.caricamentoChat = false;
      }
    },

    // 2. Apertura della chat al click
    async apriChat(conv) {
      const id = typeof conv === 'object' ? conv.id : conv;
      this.idChatAttiva = id;
      this.chatAttivaInfo = typeof conv === 'object' ? conv : this.conversazioni.find(c => c.id === id) || null;
      
      this.messaggi = [];
      this.caricamentoMessaggi = true;
      
      try {
        const res = await api.get(`/conversations/${id}/messages`);
        this.messaggi = [...(res.data || [])].reverse();
        this.$nextTick(() => this.scrollInBasso());
      } catch (e) {
        console.error("Errore nel recupero dei messaggi:", e);
        alert("Impossibile caricare i messaggi della conversazione.");
      } finally {
        this.caricamentoMessaggi = false;
      }
    },

    // 3. Invio di un nuovo messaggio (POST /conversations/{id}/messages)
    async inviaMessaggio() {
      if ((!this.nuovoMessaggio.trim() && !this.immagineSelezionata) || !this.idChatAttiva) return;

      const isImage = !!this.immagineSelezionata;
      const contenuto = isImage ? this.immagineSelezionata : this.nuovoMessaggio.trim();
      const replyId = this.messaggioRisposta ? this.messaggioRisposta.id : null;

      const payload = {
        content: contenuto,
        isImage: isImage,
        replyTo: replyId
      };

      try {
        const res = await api.post(`/conversations/${this.idChatAttiva}/messages`, payload);

        this.messaggi.push({
          id: res.data && res.data.id ? res.data.id : Date.now(),
          user: this.mioUsername,
          conversation: this.idChatAttiva,
          content: contenuto,
          timestamp: new Date().toISOString(),
          isImage: isImage,
          isForwarded: false,
          answerTo: replyId,
          state: 'sent',
          reactions: []
        });

        this.nuovoMessaggio = '';
        this.immagineSelezionata = null;
        this.nomeFileImmagine = '';
        this.messaggioRisposta = null;

        this.$nextTick(() => this.scrollInBasso());
        this.scaricaConversazioni();
      } catch (e) {
        console.error("Errore durante l'invio del messaggio:", e);
        alert("Errore nell'invio del messaggio.");
      }
    },
    selezionaFoto(event) {
      const file = event.target.files[0];
      if (!file) return;

      this.nomeFileImmagine = file.name;
      const reader = new FileReader();
      reader.onload = () => {
        this.immagineSelezionata = reader.result;
      };
      reader.readAsDataURL(file);
      event.target.value = '';
    },

    annullaImmagine() {
      this.immagineSelezionata = null;
      this.nomeFileImmagine = '';
    },

    // 4. Abbandona Gruppo (DELETE /groups/{id}/members/me)
    async abbandonaGruppo() {
      if (!this.idChatAttiva) return;
      const conferma = confirm("Sei sicuro di voler abbandonare questo gruppo?");
      if (!conferma) return;

      try {
        await api.delete(`/groups/${this.idChatAttiva}/members/me`);
        alert("Hai abbandonato il gruppo con successo.");

        this.idChatAttiva = null;
        this.chatAttivaInfo = null;
        this.messaggi = [];

        await this.scaricaConversazioni();
      } catch (err) {
        console.error("Errore nell'abbandono del gruppo:", err);
        if (err.response && err.response.data && err.response.data.error) {
          alert(err.response.data.error);
        } else {
          alert("Si è verificato un errore durante l'uscita dal gruppo.");
        }
      }
    },

    async gestisciAbbandonoDaModale(){
      this.mostraModaleInfoGruppo = false;
      await this.abbandonaGruppo();
    },

    // 5. Profilo utente loggato (GET /users?username=...)
    async caricaMioProfilo() {
      if (!this.mioUsername) return;
      try {
        const res = await api.get(`/users?username=${encodeURIComponent(this.mioUsername)}`);
        if (Array.isArray(res.data) && res.data.length > 0) {
          const mioUtente = res.data[0];
          this.miaFoto = mioUtente.proPic || '';
          if (this.miaFoto) {
            localStorage.setItem('proPic', this.miaFoto);
          }
        }
      } catch (err) {
        console.error("Errore nel recupero della foto profilo:", err);
      }
    },

    // 6. Ricerca utenti
    async attivaModalitaRicerca() {
      this.modalitaRicerca = true;
      if (this.testoRicerca.trim() === '') {
        await this.eseguiChiamataUtenti('');
      }
    },

    async cercaUtenti() {
      this.modalitaRicerca = true;
      await this.eseguiChiamataUtenti(this.testoRicerca.trim());
    },

    async eseguiChiamataUtenti(query) {
      this.caricamentoUtenti = true;
      try {
        const url = query ? `/users?username=${encodeURIComponent(query)}` : '/users';
        const res = await api.get(url);
        const lista = Array.isArray(res.data) ? res.data : [];

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

    annullaRicerca() {
      this.modalitaRicerca = false;
      this.testoRicerca = '';
      this.utentiTrovati = [];
    },

    async avviaChatDiretta(usernameSelezionato) {
      try {
        const res = await api.post('/conversations', { username: usernameSelezionato });
        this.annullaRicerca();
        await this.scaricaConversazioni();
        
        if (res.data && res.data.id) {
          const convCreata = this.conversazioni.find(c => c.id === res.data.id) || { 
            id: res.data.id, 
            name: usernameSelezionato, 
            isGroup: false 
          };
          this.apriChat(convCreata);
        }
      } catch (e) {
        alert("Errore nell'avvio della conversazione.");
        console.error(e);
      }
    },
    gestisciProfiloAggiornato(dati) {
      // Aggiorna lo stato reattivo di HomeView in tempo reale
      this.miaFoto = dati.foto;
      this.mioUsername = dati.username;

      // Ricarica la lista per allineare eventuali conversazioni o messaggi
      this.scaricaConversazioni();
    },
    gestisciGruppoCreato() {
      this.mostraModaleGruppo = false;
      this.scaricaConversazioni();
    },

    gestisciGruppoAggiornato(dati) {
    // 1. Aggiorna i metadati della chat correntemente aperta
    if (this.chatAttivaInfo && this.chatAttivaInfo.id === dati.id) {
      this.chatAttivaInfo.name = dati.name;
      this.chatAttivaInfo.picture = dati.picture;
    }

    // 2. Aggiorna l'elemento corrispondente nella lista conversazioni della sidebar
    const conv = this.conversazioni.find(c => c.id === dati.id);
    if (conv) {
      conv.name = dati.name;
      conv.picture = dati.picture;
    }
  },
    scrollInBasso() {
      const box = this.$refs.areaMessaggi;
      if (box) {
        box.scrollTop = box.scrollHeight;
      }
    },

    formattaData(dataStringa) {
      if (!dataStringa) return '';
      return new Date(dataStringa).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    },

    logout() {
      this.mostraModaleProfilo = false;
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      localStorage.removeItem('proPic');
      this.$router.push('/login');
    }
  },
  mounted() {
    this.scaricaConversazioni();
    this.caricaMioProfilo();
  },

};

</script>

<style scoped>
.dots-btn {
  font-size: 1.1rem;
  line-height: 1;
  opacity: 0.6;
}
.dots-btn:hover {
  opacity: 1;
}

.context-menu {
  position: absolute;
  top: 28px;
  z-index: 1050;
  min-width: 185px;
}

/* Messaggi inviati (a destra): si apre verso sinistra, verso l'interno */
.menu-destra {
  right: 0;
  left: auto;
}

/* Messaggi ricevuti (a sinistra): si apre verso destra, verso l'interno */
.menu-sinistra {
  left: 0;
  right: auto;
}

.emoji-btn {
  line-height: 1;
  font-size: 1.1rem;
}
.emoji-btn:hover {
  transform: scale(1.25);
}

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
  width: 400px;
  max-width: 95%;
  background-color: #ffffff !important;
}
.group-clickable {
  cursor: pointer;
  transition: background-color 0.15s ease-in-out;
}
.group-clickable:hover {
  background-color: rgba(0, 0, 0, 0.05);
}
.hover-profile:hover {
  background-color: rgba(255, 255, 255, 0.15);
}
.bg-chat {
  background-color: #f7f9fa;
}
.message-bubble {
  max-width: 75%;
  min-width: 140px;
}
.active-conversation {
  background-color: #e9ecef !important;
}
.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>