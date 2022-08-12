import axios from 'axios';

const axiosAgent = axios.create({
    baseURL: 'http://localhost:8060/api/v1',
    timeout: 10000,
    responseType: 'json',
    headers: {
        'Content-type': 'application/json'
    },
});

export type NotificationReadModel = {
    uuid: string;
    user_uuid: string;
    training_uuid: string;
    title: string;
    trainer: string;
    content: string;
    date: string;
}

const getAllNotifications = async (userUUID: string) => {
    try {
        const endpoint = `notifications/${userUUID}`;
        const response = await axiosAgent.get<NotificationReadModel[]>(endpoint);
        return (response.status !== 200) ? [] : response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
}

const deleteNotifications = async (userUUID: string) => {
    try {
        const endpoint = `notifications/${userUUID}`;
        const response = await axiosAgent.delete(endpoint);
        return response.status

    } catch (err) {
        console.error(err);
    }
}

const notificationService = {
    getAllNotifications: getAllNotifications,
    deleteNotifications: deleteNotifications,
};

export default notificationService; 