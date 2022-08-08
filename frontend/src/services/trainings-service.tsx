import {
    ParticipantAssignWriteModel,
    ParticipantGroupReadModel,
    TrainerGroupReadModel,
    TrainingGroupReadModel,
    TrainingGroupWriteModel,
    UpdateTrainigGroupWriteModel
} from './models';

import axios from 'axios';

const axiosAgent = axios.create({
    baseURL: 'http://localhost:8070/api/v1',
    timeout: 1000,
    responseType: 'json',
    headers: {
        'Content-type': 'application/json'
    },
});


const signupParticipant = async (props: ParticipantAssignWriteModel, token: string) => {
    try {
        const endpoint = `trainers/${props.trainer_uuid}/trainings/${props.trainer_group_uuid}/participants`;
        const response = await axiosAgent.post(endpoint, props.participant, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data;
    } catch (err) {
        console.error(err);
    }
};

const createTrainingGroup = async (training: TrainingGroupWriteModel, trainerUUID: string, token: string) => {
    try {
        const endpoint = `trainers/${trainerUUID}/trainings`;
        const response = await axiosAgent.post(endpoint, training, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.status
    } catch (err) {
        console.error(err);
    }
}

const getAllTrainerGroups = async (trainerUUID: string, token: string): Promise<TrainerGroupReadModel[]> => {
    try {
        const endpoint = `trainers/${trainerUUID}/trainings`;
        const response = await axiosAgent.get<TrainerGroupReadModel[]>(endpoint, {
            headers: {
                Authorization: `Bearer ${token}`
            },
        });
        return (response.status !== 200) ? [] : response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
}

const updateTrainingGroup = async (training: UpdateTrainigGroupWriteModel, trainerUUID: string, trainingUUID: string, token: string) => {
    try {
        const endpoint = `trainers/${trainerUUID}/trainings/${trainingUUID}`;
        const response = await axiosAgent.put(endpoint, training, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.status
    } catch (err) {
        console.error(err);
    }
}

const getTrainerGroup = async (trainerUUID: string, trainingUUID: string, token: string): Promise<TrainerGroupReadModel> => {
    const endpoint = `trainers/${trainerUUID}/trainings/${trainingUUID}`;
    try {
        const response = await axiosAgent.get<TrainerGroupReadModel>(endpoint, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return (response.status !== 200) ? {} as TrainerGroupReadModel : response.data;
    } catch (err) {
        return { name: "ACCESS_FORBIDDEN", description: "ACCESS_FORBIDDEN" } as TrainerGroupReadModel;
    }
}

const getAllTrainingGroups = async (): Promise<TrainingGroupReadModel[]> => {
    try {
        const response = await axiosAgent.get<TrainingGroupReadModel[]>('/trainings');
        return (response.status !== 200) ? [] : response.data;
    } catch (err) {
        console.error(err);
        return {} as TrainingGroupReadModel[];
    }
}

const getAllParticipantGroups = async (UUID: string, token: string): Promise<ParticipantGroupReadModel[]> => {
    try {
        const endpoint = `participants/${UUID}/trainings`;
        const response = await axiosAgent.get<ParticipantGroupReadModel[]>(endpoint, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return (response.status !== 200) ? [] : response.data;
    } catch (err) {
        console.error(err);
        return [];
    }
};

const cancelParticipantWorkout = async (participantUUID: string, trainerUUID: string, groupUUID: string, token: string) => {
    try {
        const endpoint = `trainers/${trainerUUID}/trainings/${groupUUID}/participants/${participantUUID}`;
        const response = await axiosAgent.delete(endpoint, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.status;
    } catch (err) {
        console.error(err)
    }
};


const deleteTrainerGroup = async (groupUUID: string, trainerUUID: string, token: string) => {
    try {
        const endpoint = `trainers/${trainerUUID}/trainings/${groupUUID}`;
        const response = await axiosAgent.delete(endpoint, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
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
    updateTrainingGroup: updateTrainingGroup,
};

export {
    TrainingsService
};
