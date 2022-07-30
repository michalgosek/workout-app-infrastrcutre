import { ENDPOINTS } from './consts';
import axios from 'axios';

export interface TrainigGroup {
    date: string;
    description: string;
    limit: number;
    name: string;
    participants: number;
    trainer_name: string;
};

const getAllTrainingGroups = async (): Promise<TrainigGroup[]> => {
    try {
        const response = await axios.get<TrainigGroup[]>(ENDPOINTS.ALL_TRAININGS);
        return response.data;
    } catch (err) {
        return [];
    }
}

const TrainingsServiceAPI = {
    GetAllTrainings: getAllTrainingGroups
};

export {
    TrainingsServiceAPI,
}; 
