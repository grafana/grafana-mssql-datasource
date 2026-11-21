import { render, screen } from '@testing-library/react';
import { openMenu } from 'react-select-event';

import { selectOptionInTest } from '../test/helpers/selectOptionInTest';

import { AzureCredentialsForm, type Props } from './AzureCredentialsForm';

const defaultProps: Props = {
  managedIdentityEnabled: false,
  workloadIdentityEnabled: false,
  azureEntraPasswordCredentialsEnabled: false,
  userIdentityEnabled: false,
  credentials: { authType: 'clientsecret' },
  onCredentialsChange: jest.fn(),
};

describe('AzureCredentialsForm', () => {
  it('does not offer workload identity when it is disabled in Grafana config', async () => {
    render(<AzureCredentialsForm {...defaultProps} />);

    openMenu(screen.getByRole('combobox'));

    expect(await screen.findByText('App Registration', { selector: '[role="option"] *' })).toBeInTheDocument();
    expect(screen.queryByText('Workload Identity')).not.toBeInTheDocument();
  });

  it('offers workload identity when it is enabled in Grafana config', async () => {
    const onCredentialsChange = jest.fn();
    render(
      <AzureCredentialsForm {...defaultProps} workloadIdentityEnabled onCredentialsChange={onCredentialsChange} />
    );

    await selectOptionInTest(screen.getByRole('combobox'), 'Workload Identity');

    expect(onCredentialsChange).toHaveBeenCalledWith({ authType: 'workloadidentity' });
  });

  it('shows workload identity as the selected auth type and no client secret fields', () => {
    render(
      <AzureCredentialsForm {...defaultProps} workloadIdentityEnabled credentials={{ authType: 'workloadidentity' }} />
    );

    expect(screen.getByText('Workload Identity')).toBeInTheDocument();
    expect(screen.queryByLabelText('Tenant ID')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('Client ID')).not.toBeInTheDocument();
    expect(screen.getByText('Configured in Grafana')).toBeInTheDocument();
  });
});
