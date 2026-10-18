<template>
  <!-- Emetto l'evento 'click' verso il padre quando l'utente preme il bottone -->
  <button
    @click="$emit('seleziona')"
    class="list-group-item list-group-item-action d-flex align-items-start p-3 border-bottom"
  >
    <!-- Foto: vera se presente, altrimenti iniziale del nome -->
    <img
      v-if="conv.picture"
      :src="conv.picture"
      class="rounded-circle me-3"
      style="width: 50px; height: 50px; object-fit: cover;"
    />
    <div 
      v-else 
      class="rounded-circle bg-primary text-white d-flex justify-content-center align-items-center me-3" 
      style="width: 50px; height: 50px; font-weight: bold;"
    >
      {{ conv.name.charAt(0).toUpperCase() }}
    </div>

    <!-- Dati conversazione -->
    <div class="flex-grow-1 overflow-hidden text-start">
      <div class="d-flex justify-content-between w-100">
        <strong class="mb-1 text-truncate">{{ conv.name }}</strong>
        <small class="text-muted">{{ formattaData(conv.lastMessageDate) }}</small>
      </div>
      <p class="mb-0 text-muted text-truncate small">
        {{ conv.lastMessagePreview || 'Nessun messaggio' }}
      </p>
    </div>
  </button>
</template>

<script>
export default {
  name: 'ConversationItem',
  // 'props' sono i dati che il padre "inietta" dentro questo componente
  props: {
    conv: {
      type: Object,
      required: true
    }
  },
  methods: {
    formattaData(dataStringa) {
      if (!dataStringa) return '';
      return new Date(dataStringa).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
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