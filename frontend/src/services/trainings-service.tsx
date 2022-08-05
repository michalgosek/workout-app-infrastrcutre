import { ParticipantAssignWriteModel, ParticipantGroupReadModel, TrainerGroupReadModel, TrainingGroupReadModel, TrainingGroupWriteModel } from './models';

import axios from 'axios';

const ENDPOINTS = {
    TRAININGS: 'http://localhost:8070/api/v1/trainings',
    TRAINERS: 'http://localhost:8070/api/v1/trainers',
    PARTICIPANTS: 'http://localhost:8070/api/v1/participants'
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

const createTrainingGroup = async (training: TrainingGroupWriteModel) => {
    try {
        const response = await axios.post(ENDPOINTS.TRAININGS, training)
        return response.status

    } catch (err) {
        console.error(err);
    }
}

const getTrainerGroup = async (trainerUUID: string, trainingUUID: string): Promise<TrainerGroupReadModel> => {
    const GET_TRAINER_GROUP_ENDPOINT = `${ENDPOINTS.TRAINERS}/${trainerUUID}/trainings/${trainingUUID}`
    try {
        const response = await axios.get<TrainerGroupReadModel>(GET_TRAINER_GROUP_ENDPOINT);
        return response.data;
    } catch (err) {
        return {} as TrainerGroupReadModel;
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



const getAllParticipantGroups = async (UUID: string): Promise<ParticipantGroupReadModel[]> => {
    try {
        const endpoint = `${ENDPOINTS.PARTICIPANTS}/${UUID}`
        const response = await axios.get<ParticipantGroupReadModel[]>(endpoint)
        return response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
};

const cancelParticipantWorkout = async (participantUUID: string, trainerUUID: string, groupUUID: string) => {
    try {
        const endpoint = `${ENDPOINTS.TRAINERS}/${trainerUUID}/trainings/${groupUUID}/participants/${participantUUID}`;
        const response = await axios.delete(endpoint)
        return response.data;
    } catch (err) {
        console.error(err)
    }
};

const getAllTrainerGroups = async (UUID: string): Promise<TrainerGroupReadModel[]> => {
    try {
        const endpoint = `${ENDPOINTS.TRAINERS}/${UUID}`
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
    getAllParticipantGroups: getAllParticipantGroups,
    cancelParticipantWorkout: cancelParticipantWorkout,
    getTrainerGroup: getTrainerGroup,
};

export {
    TrainingsService
};
