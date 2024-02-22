import {defineStore} from "pinia";

export const useMessageStore = defineStore('message', {
  state: () => ({
    messageBarText: "",
    messageBarColor: "",
    messageBarVisible: false
  }),
  actions: {
    showMessage(message: string, color: string = "primary") {
      this.messageBarText = message;
      this.messageBarColor = color;
      this.messageBarVisible = true;
    }
  }
});
