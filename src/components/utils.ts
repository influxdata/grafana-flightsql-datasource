import { useAsync } from 'react-use';
import { FlightSQLDataSource } from '../datasource';
import { SelectableValue } from '@grafana/data';

type AsyncTablesState = {
    loadingTable: boolean;
    tables: Array<SelectableValue<string>>;
    errorTable: Error | undefined;
  };
  
  type AsyncColumnsStateColumn = {
    loadingColumn: boolean;
    columns: Array<SelectableValue<string>>;
    errorColumn: Error | undefined;
  };

export const GetColumns = (datasource: FlightSQLDataSource): AsyncColumnsStateColumn => {
    const result = useAsync(async () => {
      const { columns } = await datasource.getColumns();
      return columns.map((c: any) => ({
        label: c,
        value: c,
      }));
    }, [datasource]);
  
    return {
      loadingColumn: result.loading,
      columns: result.value ?? [],
      errorColumn: result.error,
    };
  }
  
 export const GetTables = (datasource: FlightSQLDataSource): AsyncTablesState => {
    const result = useAsync(async () => {
      const res = await datasource.getTables();
      return res.frames[0].data.values[0].map((t: string) => ({
        label: t,
        value: t,
      }));
    }, [datasource]);
  
    return {
      loadingTable: result.loading,
      tables: result.value ?? [],
      errorTable: result.error,
    };
  }