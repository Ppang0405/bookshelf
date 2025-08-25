<template>
  <div class="epub-reader-container">
    <v-app-bar app color="primary" dark>
      <v-btn icon @click="goBack">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>{{ bookTitle }}</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn icon @click="prevPage" :disabled="!canGoPrev">
        <v-icon>mdi-chevron-left</v-icon>
      </v-btn>
      <v-btn icon @click="nextPage" :disabled="!canGoNext">
        <v-icon>mdi-chevron-right</v-icon>
      </v-btn>
      <v-btn icon @click="showSettings = !showSettings">
        <v-icon>mdi-cog</v-icon>
      </v-btn>
    </v-app-bar>

    <div class="epub-reader-content">
      <div id="epub-viewer" ref="epubViewer" style="height: 100%; width: 100%;"></div>
    </div>

    <v-navigation-drawer
      v-model="showSettings"
      app
      right
      temporary
      width="300"
    >
      <v-list>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title>Font Size</v-list-item-title>
            <v-slider
              v-model="fontSize"
              min="12"
              max="24"
              step="2"
              @input="updateReaderStyles"
            ></v-slider>
          </v-list-item-content>
        </v-list-item>
        
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title>Theme</v-list-item-title>
            <v-select
              v-model="theme"
              :items="themeOptions"
              @change="updateReaderStyles"
            ></v-select>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-snackbar v-model="errorSnackbar" color="error">
      {{ errorMessage }}
      <template v-slot:action="{ attrs }">
        <v-btn text v-bind="attrs" @click="errorSnackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script>
import Vue from 'vue';

export default Vue.extend({
  name: 'EpubReader',
  data() {
    return {
      epubUrl: '',
      book: null,
      rendition: null,
      location: null,
      bookTitle: 'Loading...',
      canGoPrev: false,
      canGoNext: true,
      showSettings: false,
      fontSize: 16,
      theme: 'light',
      themeOptions: [
        { text: 'Light', value: 'light' },
        { text: 'Dark', value: 'dark' },
        { text: 'Sepia', value: 'sepia' }
      ],
      readerStyles: {},
      errorSnackbar: false,
      errorMessage: ''
    };
  },
  async mounted() {
    await this.loadEpubJS();
    this.loadEpub();
    this.updateReaderStyles();
  },
  methods: {
    async loadEpubJS() {
      try {
        const epubModule = await import('epubjs');
        this.ePub = epubModule.default || epubModule;
      } catch (error) {
        console.error('Error loading epub.js:', error);
        this.showError('Failed to load EPUB reader library');
      }
    },
    
    async loadEpub() {
      try {
        const bookId = this.$route.params.bookid;
        const fileType = this.$route.params.filetype;
        
        this.epubUrl = `/api/book/${bookId}/file/${fileType}`;
        
        this.book = this.ePub(this.epubUrl);
        this.rendition = this.book.renderTo(this.$refs.epubViewer, { 
          width: '100%', 
          height: '100%' 
        });
        
        await this.rendition.display();
        
        this.book.ready.then(() => {
          this.bookTitle = this.book.metadata.title || 'Unknown Title';
        });
        
        this.rendition.on('relocated', (location) => {
          this.location = location;
          this.canGoPrev = !this.book.spine.first().href.includes(location.start.href);
          this.canGoNext = !this.book.spine.last().href.includes(location.end.href);
        });
        
      } catch (error) {
        console.error('Error loading EPUB:', error);
        this.showError('Failed to load EPUB file');
      }
    },
    
    prevPage() {
      if (this.rendition) {
        this.rendition.prev();
      }
    },
    
    nextPage() {
      if (this.rendition) {
        this.rendition.next();
      }
    },
    
    goBack() {
      this.$router.go(-1);
    },
    
    updateReaderStyles() {
      if (this.rendition) {
        const styles = {
          'font-size': `${this.fontSize}px`
        };
        
        if (this.theme === 'dark') {
          styles.color = '#ffffff';
          styles['background-color'] = '#121212';
        } else if (this.theme === 'sepia') {
          styles.color = '#5c4a37';
          styles['background-color'] = '#f4ecd8';
        } else {
          styles.color = '#000000';
          styles['background-color'] = '#ffffff';
        }
        
        this.rendition.getContents().forEach(contents => {
          contents.addStylesheet(styles);
        });
      }
    },
    
    showError(message) {
      this.errorMessage = message;
      this.errorSnackbar = true;
    }
  }
});
</script>

<style scoped>
.epub-reader-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.epub-reader-content {
  flex: 1;
  overflow: hidden;
  padding-top: 64px; /* Account for app bar */
}
</style>
