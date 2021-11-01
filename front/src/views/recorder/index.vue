<template>
<div id="app">
<div id="select">
    <label>请选择回看文件：</label>
    <select @change="selectFile($event)">
        <option v-for="(file, index) in file_list" :key="index">
            {{ file }}
        </option>
    </select>
</div>
<div id="player">{{data}}</div>
</div>
</template>
<script>
import '@/styles/asciinema-player.css'
import axios from "axios"
export default {
    name:"Rec",
    data () {
        return {
            url:'',
            file_list:[],
        }
    },
    mounted () {
        axios
        .get('/recoder')
        .then(response => {
            console.log(response)
            this.file_list = response.data
        })
        .catch(function (error) { 
            console.log(error);
        });
    },
    methods: {
        selectFile:function(event){
            console.log(event.target.value)
            this.url = "/rec/cast/"+event.target.value
            console.log(this.url)
            document.getElementById("player").innerHTML= '<asciinema-player src="'+this.url+'"></asciinema-player>'
        }
    }
}
</script>

<style scoped>
#select {
    background-color: #fff;
}
</style>