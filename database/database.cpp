#ifndef DATABASE_CPP
#define DATABASE_CPP


#include <vector>
#include <string>
#include <memory>
#include "table.cpp"

class Database {
private:
    std::vector<Table*> tables;

public:
    ~Database() {
        for (auto table : tables) {
            delete table;
        }
    }

    void addTable(Table* table) {
        tables.push_back(table);
    }

    Table* getTable(const std::string& name) {
        for (auto table : tables) {
            if (table->getName() == name) {
                return table;
            }
        }
        return nullptr;
    }

    int getTableCount() const {
        return tables.size();
    }
};

#endif