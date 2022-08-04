import React, { FC, useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
import { User, useAuth0 } from '@auth0/auth0-react';

import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import Button from 'react-bootstrap/Button';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import Form from 'react-bootstrap/Form';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { PlanTrainingGroupWriteModel } from 'services/models';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import { TrainingsService } from 'services/trainings-service';

type PlanTrainingGroupFormValues = {
    groupName: string;
    appointment: Date;
    description: string;
};

type PlanTrainingData = {
    formValues: PlanTrainingGroupFormValues;
    user: User;
};

const validateTrainingGroupData = (data: PlanTrainingData) => {
    const got: string[] = [
        data.user.name ?? '',
        data.user.sub ?? '',
        data.formValues.description,
        data.formValues.groupName,
    ]
    if (!data.formValues.appointment) return;
    return !got.some(p => p === "");
};

const createTrainingGroup = async (data: PlanTrainingData) => {
    const training: PlanTrainingGroupWriteModel = {
        user: {
            role: 'Trainer',
            uuid: data.user.sub as string,
            name: data.user.name as string,
        },
        group_name: data.formValues.groupName,
        group_desc: data.formValues.description,
        date: data.formValues.appointment.toString(),
    }

    const res = await TrainingsService.createTrainingGroup(training);
    console.log(res);
    window.location.reload();
};

const TrainingPlanningForm: FC = () => {
    const [value, setValue] = useState<Date | null>(new Date());
    const { handleSubmit, register, formState: { errors } } = useForm<PlanTrainingGroupFormValues>();
    const { user } = useAuth0();
    if (!user) {
        return null;
    }

    const onHandleChange = (newValue: Date | null) => {
        setValue(newValue);
    };
    const onSubmitHandle: SubmitHandler<PlanTrainingGroupFormValues> = (data) => {
        const trainingGroupData: PlanTrainingData = {
            formValues: {
                appointment: data.appointment,
                groupName: data.groupName,
                description: data.description
            },
            user: user,
        };
        const isDataValid = validateTrainingGroupData(trainingGroupData)
        if (!isDataValid) return;
        createTrainingGroup(trainingGroupData);
    }

    return (
        <Form onSubmit={handleSubmit(onSubmitHandle)} className='form'>
            <Form.Group className='mb-3' controlId='formBasicEmail'>
                <Form.Label>Training group name:</Form.Label>
                <Form.Control
                    type='text'
                    placeholder='Training group name'
                    {...register('groupName', { required: 'training group name is obligatory!' })}
                />
                {errors?.groupName && <p>{errors.groupName.message}</p>}
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
                            value={value}
                            onChange={onHandleChange}
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
                    placeholder='Leave a example description here'
                    {...register('description', { required: 'This value is obligatory!' })}
                    style={{ height: '100px' }}
                />
            </Form.Group>
            <Button className='btn btn-primary me-2' type='submit'>Plan</Button>
            <Button className='btn btn-secondary' type='reset' >Clear</Button>
        </Form >
    );
};

const PlanTraining: React.FC = () => {
    return (
        <>
            <TrainingPlanningForm />
        </>
    );
};

export default PlanTraining;