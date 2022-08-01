import axios from 'axios';

const ENDPOINTS = {
    ALL_TRAININGS: "http://localhost:8070/api/v1/trainings",
    PLAN_TRAINING: "http://localhost:8070/api/v1/trainings"
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


export type TrainerPOST = {
    uuid: string;
    name: string;
    role: "Trainer";
}


export type PlanTrainingGroupPOST = {
    user: TrainerPOST;
    group_name: string;
    group_desc: string;
    date: string;
}

const createTrainingGroup = async (training: PlanTrainingGroupPOST) => {
        try {
            const response = await  axios.post(ENDPOINTS.PLAN_TRAINING, training)
            return response.status

        } catch (err) {
            console.error(err);
        }
}

const getAllTrainingGroups = async (): Promise<TrainingGroup[]> => {
    try {
        const response = await axios.get<TrainingGroup[]>(ENDPOINTS.ALL_TRAININGS);
        return response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
}

const TrainingsService = {
    getAllTrainings: getAllTrainingGroups,
    createTrainingGroup: createTrainingGroup,
};

export {
    TrainingsService,
};
