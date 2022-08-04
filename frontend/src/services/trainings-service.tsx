import { ParticipantAssignWriteModel, PlanTrainingGroupWriteModel, TrainerGroupReadModel, TrainingGroupReadModel } from './models';

import axios from 'axios';

const ENDPOINTS = {
    TRAININGS: 'http://localhost:8070/api/v1/trainings',
    TRAINERS: 'http://localhost:8070/api/v1/trainers'
};


const signupParticipant = async (props: ParticipantAssignWriteModel) => {
    try {
        const endpoint = `${ENDPOINTS.TRAINERS}/${props.trainer_uuid}/trainings/${props.trainer_group_uuid}/participants`
        const response = await axios.post(endpoint, props.participant)
        return response.data;
    } catch (err) {
        console.error(err);
    }
};

const createTrainingGroup = async (training: PlanTrainingGroupWriteModel) => {
    try {
        const response = await axios.post(ENDPOINTS.TRAININGS, training)
        return response.status

    } catch (err) {
        console.error(err);
    }
}

const getAllTrainingGroups = async (): Promise<TrainingGroupReadModel[]> => {
    try {
        const response = await axios.get<TrainingGroupReadModel[]>(ENDPOINTS.TRAININGS);
        return response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
}

const getAllTrainerGroups = async (groupUUID: string): Promise<TrainerGroupReadModel[]> => {
    try {
        const endpoint = `${ENDPOINTS.TRAINERS}/${groupUUID}`
        const response = await axios.get<TrainerGroupReadModel[]>(endpoint)
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
    signupParticipant: signupParticipant,
};

export {
    TrainingsService
};
