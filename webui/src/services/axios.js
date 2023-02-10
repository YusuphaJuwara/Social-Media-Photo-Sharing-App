import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
	headers: {
    Authorization: "Bearer " + localStorage.getItem('token')
	}
});

export default instance;
