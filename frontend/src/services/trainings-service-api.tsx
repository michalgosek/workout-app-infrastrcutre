import axios from 'axios';

const ENDPOINTS = {
    ALL_TRAININGS: "http://localhost:8070/api/v1/trainings"
}

export interface TrainingGroup {
    uuid: string
    date: string;
    description: string;
    limit: number;
    name: string;
    participants: number;
    trainer_name: string;
};

const getAllTrainingGroups = async (): Promise<TrainingGroup[]> => {
    try {
        const response = await axios.get<TrainingGroup[]>(ENDPOINTS.ALL_TRAININGS);
        return response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
}

const TrainingsServiceAPI = {
    getAllTrainings: getAllTrainingGroups
};

export {
    TrainingsServiceAPI,
};
