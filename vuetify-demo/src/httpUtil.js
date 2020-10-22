import axios from 'axios'

export default {
    name: 'httpUtil',
    uri: {
        chat: "http://localhost:8081/api/v1/chat",
        usr: "http://localhost:8081/api/v1/usr",
    },
    post(uri, args, succ, fail) {
        axios.post(uri, JSON.stringify(args), { withCredentials: true }).then((response) => {
            console.log("%s\n%s", uri, JSON.stringify(response.data));
            succ(response.data, args);
        }).catch((error) => {
            console.log("%s\n%s", uri, JSON.stringify(error))
            fail(error, args);
        })
    },
    get(uri, args, succ, fail) {
        axios.get(uri, { params: args, withCredentials: true }).then(function (response) {
            console.log("%s\n%s", uri, JSON.stringify(response.data));
            succ(response.data, args);
        }).catch(function (error) {
            console.log("%s\n%s", uri, JSON.stringify(error))
            fail(error, args);
        });
    },
}
