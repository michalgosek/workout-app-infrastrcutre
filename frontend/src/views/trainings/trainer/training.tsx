import { Container, Table } from 'react-bootstrap';
import { FC, FormEventHandler, PropsWithChildren, useEffect, useState } from 'react';
import { ParticipantReadModel, TrainerGroupReadModel } from 'services/models';
import { SubmitHandler, useForm } from 'react-hook-form';

import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import Button from 'react-bootstrap/Button';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import DropButton from 'components/buttons/drop-button';
import Form from 'react-bootstrap/Form';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import { TrainingsService } from 'services/trainings-service';
import style from './training.module.scss';
import { useParams } from 'react-router-dom';

type ParticipantsTableProps = {
    trainingUUID: string;
    trainerUUID: string;
    participants: ParticipantReadModel[];
};

type deleteTrainingData = {
    trainerUUID?: string;
    particiapntUUID?: string;
    trainingUUID?: string;
};

const deleteParticipant = (data: deleteTrainingData) => {
    console.log(data);
};

const ParticipantsTable: FC<PropsWithChildren<ParticipantsTableProps>> = ({ participants, trainerUUID, trainingUUID }) => {
    return (
        <div className={style.rowtext}>
            <Table striped bordered hover>
                <thead className={style['t-headers']}>
                    <tr>
                        <th>#</th>
                        <th>Participant Name</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    {
                        participants.map((p, idx) => (
                            <tr key={p.uuid} className={style['t-rows']} >
                                <td>{idx + 1}</td>
                                <td>{p.name}</td>
                                <td>
                                    <DropButton onClickHandle={() => deleteParticipant({
                                        particiapntUUID: p.uuid,
                                        trainerUUID: trainerUUID,
                                        trainingUUID: trainingUUID,
                                    })} />
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );
};

const NoTrainerGroupParticipants: FC = () => {
    return <h1>No participants.</h1>
};

const NoTrainerGroup: FC = () => {
    return <h1>Trainer Group Not Exist.</h1>
};



type PlanTrainingGroupFormValues = {
    groupName: string;
    appointment: Date;
    description: string;
};



const TrainerGroup: FC = () => {
    const { trainingUUID, trainerUUID } = useParams();
    const [training, setTraining] = useState<TrainerGroupReadModel>();
    const [value, setValue] = useState<Date | null>(new Date());
    const { handleSubmit, register, formState: { errors } } = useForm<PlanTrainingGroupFormValues>();



    useEffect(() => {
        if (!trainerUUID || !trainingUUID) {
            console.error('missing trainerUUID or trainingUUID in URL path');
            return;
        }
        TrainingsService.getTrainerGroup(trainerUUID, trainingUUID)
            .then(t => {
                setTraining(t);
                setValue(new Date(t.date));
            })
            .catch(err => console.error(err));
    }, [trainerUUID, trainingUUID]);

    if (!training || !trainerUUID || !trainingUUID) return <NoTrainerGroup />;

    const onSubmitHandle: SubmitHandler<TrainingFormValues> = (data) => {
        console.log(data);
    };

    return (
        <>
            <Container>
                <TrainingForm handleSubmitCallback={onSubmitHandle} placeholders={{
                    name: training.name,
                    description: training.description,
                    date: new Date(training.date)
                }} />
            </Container>

            <Container className='mt-3'>
                {!training.participants ? <NoTrainerGroupParticipants /> : <ParticipantsTable participants={training.participants} trainerUUID={trainerUUID} trainingUUID={trainingUUID} />}
            </Container>
        </>
    );
};

type TrainingFormValues = {
    description: string;
    name: string;
    date: Date;
}

type TrainingFormProps = {
    handleSubmitCallback: SubmitHandler<TrainingFormValues>
    placeholders: TrainingFormValues;
};


// nie zmienia się data, a wszystko inne tak!

const TrainingForm: FC<PropsWithChildren<TrainingFormProps>> = ({ handleSubmitCallback, placeholders }) => {
    const [appointment, setAppointment] = useState<Date | null>(new Date());
    const { handleSubmit, register, formState: { errors } } = useForm<TrainingFormValues>();
    const onHandleAppointment = (updated: Date | null) => {
        // console.log(updated);
        setAppointment(updated);
    };

    return (
        <Form onSubmit={handleSubmit(handleSubmitCallback)} className='form'>
            <Form.Group className='mb-3' controlId='formBasicEmail'>
                <Form.Label>Training group name:</Form.Label>
                <Form.Control
                    type='text'
                    placeholder={placeholders.name}
                    {...register('name', { required: 'training group name is obligatory!' })}
                />
                {errors?.name && <p>{errors.name.message}</p>}
                <Form.Text className='text-muted'>
                    This will be used as group name description.
                </Form.Text>
            </Form.Group>
            <Form.Group className='mb-3' controlId='formBasicDate'>
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                    <Stack spacing={3}>
                        <DateTimePicker
                            {...register('date', { required: 'appointment date for training group is obligatory!' })}
                            label='Training appointment'
                            value={appointment}
                            onChange={onHandleAppointment}
                            ampm={false}
                            renderInput={(params) => <TextField {...params} />}
                        />
                    </Stack>
                </LocalizationProvider>
            </Form.Group>
            <Form.Group className='mb-3' controlId='formBasicTextArea'>
                <Form.Label>Training group description:</Form.Label>
                <Form.Control
                    as='textarea'
                    placeholder={placeholders.description}
                    {...register('description', { required: 'This value is obligatory!' })}
                    style={{ height: '100px' }}
                />
            </Form.Group>
            <Button className='btn btn-primary me-2' type='submit'>Plan</Button>
            <Button className='btn btn-secondary' type='reset' >Clear</Button>
        </Form >
    );
}

export default TrainerGroup;