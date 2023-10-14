import axios from 'axios';

const baseURL = window.location.protocol + "//" +
    window.location.hostname +
    (window.location.port ? ":" + window.location.port : "") + "/";

const api = axios.create({
    baseURL: baseURL,
});

export default api;