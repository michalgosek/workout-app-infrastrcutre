import axios from 'axios';

const ENDPOINTS = {
    TRAININGS: "http://localhost:8070/api/v1/trainings",
    TRAINERS: "http://localhost:8070/api/v1/trainers"
};

export type TrainingGroup = {
    uuid: string
    date: string;
    description: string;
    limit: number;
    name: string;
    participants: number;
    trainer_name: string;
};


type Trainer = {
    uuid: string;
    name: string;
    role: "Trainer";
}


export type PlanTrainingGroupPOST = {
    user: Trainer;
    group_name: string;
    group_desc: string;
    date: string;
}

const createTrainingGroup = async (training: PlanTrainingGroupPOST) => {
    try {
        const response = await axios.post(ENDPOINTS.TRAININGS, training)
        return response.status

    } catch (err) {
        console.error(err);
    }
}


type Participant = {
    uuid: string;
    name: string;
}

export type TrainerGroup = {
    uuid: string;
    name: string;
    description: string;
    date: string;
    limit: number;
    participant: Participant[];
}


const getAllTrainingGroups = async (): Promise<TrainingGroup[]> => {
    try {
        const response = await axios.get<TrainingGroup[]>(ENDPOINTS.TRAININGS);
        return response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
}

const getAllTrainerGroups = async (groupUUID: string): Promise<TrainerGroup[]> => {
    try {
        const endpoint = `${ENDPOINTS.TRAINERS}/${groupUUID}`
        const response = await axios.get<TrainerGroup[]>(endpoint)
        return response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
}

const deleteTrainerGroup = async (groupUUID: string, trainerUUID: string) => {
    try {
        const endpoint = `http://localhost:8070/api/v1/trainers/${trainerUUID}/trainings/${groupUUID}`
        const response = await axios.delete(endpoint)
        return response.data;
    } catch (err) {
        console.error(err);
    }
}




const TrainingsService = {
    getAllTrainings: getAllTrainingGroups,
    createTrainingGroup: createTrainingGroup,
    getAllTrainerGroups: getAllTrainerGroups,
    deleteTrainerGroup: deleteTrainerGroup,
};

export {
    TrainingsService,
};
