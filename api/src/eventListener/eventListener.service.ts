import {
  BadRequestException,
  Injectable,
  NotFoundException,
  UnauthorizedException,
} from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { SupportedChainEntity } from "common/database/entities/supportedChain.entity";
import { ClientEntity } from "common/database/entities/client.entity";
import { EventListenerEntity } from "common/database/entities/eventListener.entity";
import { CreateEventListenerDto } from "./dto/create-eventListener.dto";
import { CreateEventListenerResultDto } from "./dto/create-eventListener.result.dto";
import { GetEventListenerResultDto } from "./dto/get-eventListener.result.dto";
import { UUID } from "crypto";

@Injectable()
export class EventListenerService {
  constructor(
    @InjectRepository(EventListenerEntity)
    private eventListenerRepository: Repository<EventListenerEntity>,
    @InjectRepository(SupportedChainEntity)
    private supportedChainRepository: Repository<SupportedChainEntity>,
    @InjectRepository(ClientEntity)
    private clientRepository: Repository<ClientEntity>
  ) {}

  async create(createEventListenerDto: CreateEventListenerDto, clientId: UUID) {
    const supportedChain = await this.supportedChainRepository.findOneBy({
      chain: createEventListenerDto.chain,
    });

    if (!supportedChain)
      throw new BadRequestException(
        `EventListenerService: chain ${createEventListenerDto.chain} is not supported`
      );

    const eventListener = await this.eventListenerRepository.create({
      ...createEventListenerDto,
      createdTimestamp: new Date(),
      clientId: clientId,
      active: true,
    });

    await this.eventListenerRepository.save(eventListener);

    return {
      eventListenerId: eventListener.id,
    } as CreateEventListenerResultDto;
  }

  async stop(eventListenerId: UUID, clientId: UUID) {
    const eventListener = await this.eventListenerRepository.findOne({
      where: {
        id: eventListenerId,
        clientId: clientId,
      },
    });

    if (!eventListener)
      throw new NotFoundException(
        `EventListener id=${eventListenerId} not found`
      );

    eventListener.active = false;
    await this.eventListenerRepository.save(eventListener);
  }

  async findOne(id: UUID, clientId: UUID): Promise<GetEventListenerResultDto> {
    const eventListener = await this.eventListenerRepository.findOne({
      where: {
        id: id,
        clientId: clientId,
      },
    });

    if (!eventListener)
      throw new NotFoundException(`EventListener id=${id} not found`);
    return eventListener as GetEventListenerResultDto;
  }

  async findAll(clientId: UUID): Promise<GetEventListenerResultDto[]> {
    const eventListeners = await this.eventListenerRepository.find({
      where: {
        clientId: clientId,
      },
    });

    return eventListeners as GetEventListenerResultDto[];
  }
}
