import { FC, PropsWithChildren } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import { TrainingGroupWriteModel, UpdateTrainigGroupWriteModel } from 'services/models';

import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import Button from 'react-bootstrap/Button';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import Form from 'react-bootstrap/Form';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import { useAuth0 } from '@auth0/auth0-react';
import { useGetTrainingsServiceAccessToken } from 'services/hooks';
import { useParams } from 'react-router-dom';
import { useState } from 'react';

export type TrainingFormValues = {
    description: string;
    name: string;
    appointment?: string;
}

export type TrainingFormProps = {
    placeholders: TrainingFormValues;
    callbackPostAPI?: (data: TrainingGroupWriteModel, trainerUUID: string, token: string) => void;
    callbackPutAPI?: (data: UpdateTrainigGroupWriteModel, trainerUUID: string, trainingUUID: string, token: string) => void;
};

const praseToTimeFormat = (appointment: Date) => {
    const date = appointment.toLocaleDateString();
    let timeString = appointment.toLocaleTimeString();
    const time = timeString.substring(0, timeString.length - 3);
    const format = `${date} ${time}`;
    return format;
}

const parseStringToDate = (appointment: string) => {
    const datePart = appointment.slice(0, 10);
    const [hour, min] = appointment.slice(11, appointment.length).split(':');
    const [month, day, year] = datePart.split('/');
    const date = new Date(+year, +month - 1, +day);
    date.setHours(+hour);
    date.setMinutes(+min)
    return date;
}

const TrainingForm: FC<PropsWithChildren<TrainingFormProps>> = ({ placeholders, callbackPostAPI, callbackPutAPI }) => {
    const { trainerUUID, trainingUUID } = useParams();
    const [appointment, setAppointment] = useState<Date | null>(placeholders.appointment ? parseStringToDate(placeholders.appointment) : new Date());
    const { handleSubmit, register, formState: { errors } } = useForm<TrainingFormValues>();
    const { user } = useAuth0();
    const token = useGetTrainingsServiceAccessToken();

    const onHandleAppointmentChange = (update: Date | null) => {
        setAppointment(update);
    };
    if (!user) {
        console.error('null user value');
        return null;
    }
    if (callbackPutAPI && callbackPostAPI) {
        console.error("only single API must be provided for training group service: POST or PUT")
        return null;
    }
    if (!callbackPutAPI && !callbackPostAPI) {
        console.error("at least one callback API must be provided for training group service: POST or PUT")
        return null;
    }


    const TRAININGS_SERVICE_PUT_API_CALLBACK = callbackPutAPI;
    const TRAININGS_SERVICE_POST_API_CALLBACK = callbackPostAPI;

    const onSubmitHandle: SubmitHandler<TrainingFormValues> = (data) => {
        if (!appointment) {
            console.error('null appointment value');
            return;
        }
        if (!trainerUUID) {
            console.error('missing trainer UUID for TrainingsServcie PUT/POST API request');
            return;
        }

        data.appointment = praseToTimeFormat(appointment);
        if (TRAININGS_SERVICE_POST_API_CALLBACK) {
            const training: TrainingGroupWriteModel = {
                user: {
                    role: 'Trainer',
                    uuid: user.sub ?? '',
                    name: user.name ?? '',
                },
                group_name: data.name,
                group_desc: data.description,
                date: data.appointment,
            }
            TRAININGS_SERVICE_POST_API_CALLBACK(training, trainerUUID, token);
            window.location.reload();
            return;
        }

        if (TRAININGS_SERVICE_PUT_API_CALLBACK) {
            if (!trainingUUID) {
                console.error('missing training UUID for TrainingsServcie PUT API request');
                return;
            }
            const training: UpdateTrainigGroupWriteModel = {
                group_name: data.name,
                group_desc: data.description,
                date: data.appointment,
            }
            TRAININGS_SERVICE_PUT_API_CALLBACK(training, trainerUUID, trainingUUID, token);
            window.location.reload();
            return;
        }
    };

    return (
        <Form onSubmit={handleSubmit(onSubmitHandle)} className='form'>
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
                            {...register('appointment', { required: 'appointment date for training group is obligatory!' })}
                            label='Training appointment'
                            value={appointment}
                            onChange={onHandleAppointmentChange}
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

export default TrainingForm;