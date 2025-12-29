# Refatorar Autenticação para Backend

## Objetivo

Mover login/register do frontend (cliente Supabase direto) para API do backend.

## Plano de Ação

### 1. Backend - Melhorar Handlers de Autenticação

- [x] Implementar `SignIn` e `SignUp` como handlers globais
- [x] Garantir que retornem tokens JWT válidos
- [x] Adicionar validação de entrada de dados

### 2. Frontend - Remover Cliente Supabase

- [x] Modificar `auth/services.ts` para usar API backend
- [x] Remover dependência do `supabaseClient.ts`
- [x] Atualizar tipos para resposta da API

### 3. Frontend - Atualizar Store/Auth

- [ ] Modificar actions para chamar API backend
- [ ] Remover referências ao cliente Supabase
- [ ] Ajustar tratamento de erros

### 4. Testes

- [ ] Testar login via API backend
- [ ] Testar register via API backend
- [ ] Verificar fluxo completo de autenticação

### 5. Limpeza

- [x] Remover `supabaseClient.ts`
- [ ] Remover dependências Supabase do frontend (se não usadas)
- [ ] Atualizar documentação

## Arquivos a Modificar

### Backend

- `/internal/handlers/auth.go` - Adicionar handlers globais
- `/internal/routes/routes.go` - Já está correto

### Frontend

- `/src/store/auth/services.ts` - Mudar para API backend
- `/src/store/auth/actions.ts` - Atualizar chamadas
- `/src/lib/supabaseClient.ts` - Possivelmente remover

## Considerações

- Perderemos cache e gerenciamento de sessão do cliente Supabase
- Teremos que implementar refresh de tokens manualmente se necessário
- Backend terá controle total da autenticação
