
export type TrainingGroupReadModel = {
    uuid: string
    date: string;
    description: string;
    limit: number;
    name: string;
    participants: number;
    trainer_name: string;
    trainer_uuid: string;
};


type TrainerWriteModel = {
    uuid: string;
    name: string;
    role: 'Trainer';
};

export type PlanTrainingGroupWriteModel = {
    user: TrainerWriteModel;
    group_name: string;
    group_desc: string;
    date: string;
};

type ParticipantReadModel = {
    uuid: string;
    name: string;
};


export type TrainerGroupReadModel = {
    uuid: string;
    name: string;
    description: string;
    date: string;
    limit: number;
    participants: ParticipantReadModel[];
};


export type ParticipantWriteModel = {
    participant_name: string;
    participant_uuid: string;
};

export type ParticipantAssignWriteModel = {
    trainer_uuid: string;
    trainer_group_uuid: string;
    participant: ParticipantWriteModel;
}
